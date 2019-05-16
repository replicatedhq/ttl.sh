# replreg hooks

This is the hooks server that also has the reaper for replreg.is.

The webhooks from the registry are terminated in [src/controllers/HookAPI.ts](https://github.com/replicatedhq/replreg/blob/master/hooks/src/controllers/HookAPI.ts). The reaper runs on a cron and executes [src/commands/reap.ts](https://github.com/replicatedhq/replreg/blob/master/hooks/src/commands/reap.ts).



