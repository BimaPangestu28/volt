// MongoDB initialization script
// This script runs when the MongoDB container starts for the first time

db = db.getSiblingDB('volt');

// Create collections with validation
db.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['username', 'email', 'password_hash'],
      properties: {
        username: {
          bsonType: 'string',
          minLength: 3,
          maxLength: 50
        },
        email: {
          bsonType: 'string',
          pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$'
        },
        password_hash: {
          bsonType: 'string'
        },
        created_at: {
          bsonType: 'date'
        },
        updated_at: {
          bsonType: 'date'
        }
      }
    }
  }
});

db.createCollection('workspaces');
db.createCollection('collections');
db.createCollection('requests');
db.createCollection('environments');
db.createCollection('websocket_connections');
db.createCollection('comments');
db.createCollection('execution_history');

// Create indexes
db.users.createIndex({ 'email': 1 }, { unique: true });
db.users.createIndex({ 'username': 1 }, { unique: true });

db.workspaces.createIndex({ 'owner_id': 1 });
db.workspaces.createIndex({ 'members': 1 });

db.collections.createIndex({ 'workspace_id': 1 });
db.collections.createIndex({ 'created_by': 1 });

db.requests.createIndex({ 'collection_id': 1 });
db.requests.createIndex({ 'created_by': 1 });

db.environments.createIndex({ 'workspace_id': 1 });

db.websocket_connections.createIndex({ 'user_id': 1 });
db.websocket_connections.createIndex({ 'workspace_id': 1 });

db.comments.createIndex({ 'request_id': 1 });
db.comments.createIndex({ 'user_id': 1 });

db.execution_history.createIndex({ 'request_id': 1 });
db.execution_history.createIndex({ 'user_id': 1 });
db.execution_history.createIndex({ 'executed_at': -1 });

print('Volt database initialized successfully!');
print('Collections created with proper indexes.');
print('Ready for application startup.');