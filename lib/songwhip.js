import chalk from 'chalk';
import cheerio from 'cheerio';
import clipboardy from 'clipboardy';
import fetch from 'node-fetch';

export const availableServices = [
  'spotify',
  'itunes',
  'youtube',
  'tidal',
  'amazonMusic',
  'pandora',
  'deezer',
  'audiomack',
  'qobuz'
]

export function listServices() {
  availableServices.forEach((service) => {
    console.log(chalk.yellowBright(service));
  });
}


export function getSongwhipUrl(url, opts) {
  const service = opts.s ?? null;
  console.log(chalk.yellowBright(`Getting ${service || 'Songwhip'} link for ${url}`));

  fetch('https://songwhip.com/', {
    method: 'POST',
    body: `{ "url": ${JSON.stringify(url)} }`,
    headers: { 'Content-Type': 'application/json' },
  })
  .then((response) => response.json())
  .then((payload) => crawlSongwhipPage(service, payload.url))
  .then((url) => {
    clipboardy.writeSync(url);
    console.log(chalk.green(`Copied Songwhip link ${url} to clipboard!`));
  });
}

function crawlSongwhipPage(service, url) {
  if (!service) {
    return Promise.resolve(url);
  }

  if (!availableServices.includes(service)) {
    return Promise.reject('Invalid service!');
  }

  return fetch(url)
    .then((response) => response.text())
    .then((html) => {
      const dom = cheerio.load(html);
      const link = dom(`a[data-testid="ServiceButton ${service} itemLinkButton ${service}ItemLinkButton"]`);
      return link.attr('href');
    });
}
