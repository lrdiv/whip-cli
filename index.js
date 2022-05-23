#!/usr/bin/env node

import { program } from 'commander';
import { getSongwhipUrl } from './lib/songwhip.js';
import { getServiceUrl, listServices } from './lib/service.js';

program.command('song <track>')
  .description('Get a Songwhip page link from a given track')
  .action(getSongwhipUrl);

program.command('services')
  .description('List available services for URLs')
  .action(listServices);

program.command('get <service> <track>')
  .description('Get a link from a given service for a given track')
  .action(getServiceUrl);

program.parse();
