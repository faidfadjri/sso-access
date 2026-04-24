module.exports = {
  apps: [
    {
      name: "access-next",
      script: "node_modules/next/dist/bin/next",
      args: "start -p 4001",
      instances: "max",
      exec_mode: "cluster",
      watch: false,
      autorestart: true,
      max_memory_restart: "500M",
      env: {
        NODE_ENV: "production",
        NEXT_PUBLIC_DEBUG: "false"
      }
    }
  ]
};