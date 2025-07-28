#!/usr/bin/env python3
"""
Volt Deployment Webhook Server
Simple HTTP server untuk menerima deployment trigger dari GitHub Actions
"""

import json
import os
import subprocess
import threading
from http.server import BaseHTTPRequestHandler, HTTPServer
import logging

# Configuration
PORT = int(os.getenv('WEBHOOK_PORT', 9000))
TOKEN = os.getenv('WEBHOOK_TOKEN', 'your-secret-webhook-token')
LOG_LEVEL = os.getenv('LOG_LEVEL', 'INFO')

# Setup logging
logging.basicConfig(
    level=getattr(logging, LOG_LEVEL.upper()),
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class DeploymentHandler(BaseHTTPRequestHandler):
    def log_message(self, format, *args):
        """Override to use our logger"""
        logger.info(format % args)
    
    def do_POST(self):
        if self.path == '/deploy':
            self.handle_deployment()
        else:
            self.send_response(404)
            self.end_headers()
            self.wfile.write(b'{"error": "Not found"}')
    
    def do_GET(self):
        if self.path == '/health':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(b'{"status": "healthy", "service": "volt-deployment-webhook"}')
        else:
            self.send_response(404)
            self.end_headers()
            self.wfile.write(b'{"error": "Not found"}')
    
    def handle_deployment(self):
        try:
            # Get request body
            content_length = int(self.headers.get('Content-Length', 0))
            post_data = self.rfile.read(content_length)
            
            # Check authorization
            auth_header = self.headers.get('Authorization', '')
            if not auth_header.startswith('Bearer '):
                logger.warning("Missing or invalid authorization header")
                self.send_error(401, "Unauthorized")
                return
            
            token = auth_header[7:]  # Remove 'Bearer '
            if token != TOKEN:
                logger.warning("Invalid token provided")
                self.send_error(401, "Invalid token")
                return
            
            # Parse JSON payload
            try:
                payload = json.loads(post_data.decode('utf-8'))
                logger.info(f"Received deployment payload: {payload}")
            except json.JSONDecodeError:
                logger.error("Invalid JSON payload")
                self.send_error(400, "Invalid JSON")
                return
            
            # Validate payload
            if not self.validate_payload(payload):
                self.send_error(400, "Invalid payload format")
                return
            
            # Send immediate response
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(b'{"status": "deployment started"}')
            
            # Start deployment in background thread
            thread = threading.Thread(
                target=self.deploy_volt,
                args=(payload['environment'], payload['images']['backend'], payload['images']['frontend'])
            )
            thread.daemon = True
            thread.start()
            
        except Exception as e:
            logger.error(f"Error handling deployment: {e}")
            self.send_error(500, f"Internal server error: {str(e)}")
    
    def validate_payload(self, payload):
        """Validate deployment payload format"""
        required_fields = ['action', 'environment', 'images']
        for field in required_fields:
            if field not in payload:
                logger.error(f"Missing required field: {field}")
                return False
        
        if 'backend' not in payload['images'] or 'frontend' not in payload['images']:
            logger.error("Missing backend or frontend image")
            return False
        
        if payload['environment'] not in ['production', 'staging']:
            logger.error(f"Invalid environment: {payload['environment']}")
            return False
        
        return True
    
    def deploy_volt(self, environment, backend_image, frontend_image):
        """Deploy Volt application"""
        try:
            logger.info(f"🚀 Starting deployment to {environment}")
            logger.info(f"Backend: {backend_image}")
            logger.info(f"Frontend: {frontend_image}")
            
            # Change to volt directory
            os.chdir('/opt/volt')
            
            if environment == 'production':
                self.deploy_production(backend_image, frontend_image)
            elif environment == 'staging':
                self.deploy_staging(backend_image, frontend_image)
            
            logger.info(f"✅ {environment.title()} deployment completed!")
            
        except Exception as e:
            logger.error(f"❌ Deployment failed: {e}")
    
    def deploy_production(self, backend_image, frontend_image):
        """Deploy to production environment"""
        logger.info("Deploying to production...")
        
        # Pull images
        self.run_command(f"docker pull {backend_image}")
        self.run_command(f"docker pull {frontend_image}")
        
        # Update services
        backend_result = self.run_command(
            f"docker service update --image {backend_image} volt-prod_backend",
            check=False
        )
        
        if backend_result.returncode != 0:
            logger.info("Backend service not found, deploying stack...")
            self.run_command("docker stack deploy -c docker-compose.prod.yml volt-prod")
        
        frontend_result = self.run_command(
            f"docker service update --image {frontend_image} volt-prod_frontend",
            check=False
        )
        
        if frontend_result.returncode != 0:
            logger.info("Frontend service not found, deploying stack...")
            self.run_command("docker stack deploy -c docker-compose.prod.yml volt-prod")
        
        # Cleanup
        self.run_command("docker image prune -f", check=False)
    
    def deploy_staging(self, backend_image, frontend_image):
        """Deploy to staging environment"""
        logger.info("Deploying to staging...")
        
        # Pull images
        self.run_command(f"docker pull {backend_image}")
        self.run_command(f"docker pull {frontend_image}")
        
        # Update services
        backend_result = self.run_command(
            f"docker service update --image {backend_image} volt-staging_backend",
            check=False
        )
        
        if backend_result.returncode != 0:
            logger.info("Backend service not found, deploying stack...")
            self.run_command("docker stack deploy -c docker-compose.staging.yml volt-staging")
        
        frontend_result = self.run_command(
            f"docker service update --image {frontend_image} volt-staging_frontend",
            check=False
        )
        
        if frontend_result.returncode != 0:
            logger.info("Frontend service not found, deploying stack...")
            self.run_command("docker stack deploy -c docker-compose.staging.yml volt-staging")
    
    def run_command(self, command, check=True):
        """Run shell command and log output"""
        logger.info(f"Running: {command}")
        
        result = subprocess.run(
            command,
            shell=True,
            capture_output=True,
            text=True
        )
        
        if result.stdout:
            logger.info(f"STDOUT: {result.stdout.strip()}")
        
        if result.stderr:
            logger.warning(f"STDERR: {result.stderr.strip()}")
        
        if check and result.returncode != 0:
            raise subprocess.CalledProcessError(result.returncode, command)
        
        return result

def main():
    """Start the webhook server"""
    server_address = ('', PORT)
    httpd = HTTPServer(server_address, DeploymentHandler)
    
    logger.info(f"🚀 Starting Volt Deployment Webhook Server on port {PORT}")
    logger.info(f"📡 Webhook URL: http://your-server:{PORT}/deploy")
    logger.info(f"🏥 Health check: http://your-server:{PORT}/health")
    logger.info(f"🔑 Token: {TOKEN}")
    
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        logger.info("👋 Shutting down webhook server...")
        httpd.shutdown()

if __name__ == '__main__':
    main()