import chalk from 'chalk';
import cheerio from 'cheerio';
import clipboardy from 'clipboardy';
import fetch from 'node-fetch';

const availableServices = [
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

export function getServiceUrl(service, url) {
  console.log(chalk.yellowBright(`Getting ${service} link for ${url}`));

  fetch('https://songwhip.com/', {
    method: 'POST',
    body: `{ "url": ${JSON.stringify(url)} }`,
    headers: { 'Content-Type': 'application/json' },
  })
  .then((response) => response.json())
  .then((payload) => crawlSongwhipPage(service, payload.url))
  .then((serviceUrl) => {
    clipboardy.writeSync(serviceUrl);
    console.log(chalk.green(`Copied ${serviceUrl} to clipboard!`));
  })
  .catch((e) => console.log(chalk.red(e)));
}

function crawlSongwhipPage(service, url) {
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
