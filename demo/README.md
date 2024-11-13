## Start the failover playground
```bash
K8GB_LOCAL_VERSION=test FULL_LOCAL_SETUP_WITH_APPS=true make deploy-full-local-setup
```

## Start the demo UI
Navigate to the demo folder
```bash
cd demo
```

## Download dependencies
```bash
npm install express cors
```

## Start the backend server:
```bash
node server.js
```

## Start the frontend:
```bash
live-server
```
