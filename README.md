# whip CLI

Gets you a Songwhip URL for a track OR get a specific service URL

## Installation

**Clone this repository:**

```bash
git clone git@github.com:lrdiv/whip-cli.git
```

**Install the CLI:**

```bash
cd whip-cli
npm i -g
```

## Usage

The last argument for these commands can be a URL from any available streaming service (see below for how to list available services)

**Get a Songwhip URL:**

```bash
whip song https://tidal.com/browse/track/177186841
```

**Get a specific service URL:**

```bash
whip get spotify https://tidal.com/browse/track/177186841
```

The Songwhip or service URL will be copied to your clipboard for easy sharing! Wow!

**To list all available services:**

```bash
whip services
```

