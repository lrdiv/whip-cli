#!/usr/bin/env node

import { program } from 'commander';
import { getSongwhipUrl, listServices } from './lib/songwhip.js';

program.command('services')
  .description('List available services for URLs')
  .action(listServices);

program.command('get <track>')
  .option('-s <service>')
  .description('Get a shareable link from a given track')
  .action(getSongwhipUrl);

program.parse();
