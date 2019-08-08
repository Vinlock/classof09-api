const AWS = require('aws-sdk');
const fs = require('fs');
const minimist = require('minimist');
const ora = require('ora');

// Get arguments
const args = minimist(process.argv.slice(2), {
    string: [
        'secret',
        'file',
        'region',
        'aws-key',
        'aws-secret',
    ],
    alias: {
        s: 'secret',
        f: 'file',
        r: 'region',
    },
    default: {
        region: 'us-west-2',
    },
});
const secret = args['secret'];
const ymlFile = args['file'];
const region = args['region'];
const accessKeyId = args['aws-key'];
const secretAccessKey = args['aws-secret'];

// Validate
if (!secret || secret.length === 0) {
    console.error('ERROR: --secret|-s must be provided');
    process.exit(1);
}
if (!ymlFile || ymlFile.length === 0) {
    console.error('ERROR: --file|-f must be provided');
    process.exit(1);
}
if (!ymlFile.endsWith('.yml') && !ymlFile.endsWith('.yaml')) {
    console.error('ERROR: --file|-f must be a YAML file.');
    process.exit(1);
}

// Instantiate Secrets Manager
const options = {
    apiVersion: '2017-10-17',
    region,
};
if (accessKeyId && accessKeyId.length > 1) {
    options.accessKeyId = accessKeyId;
}
if (secretAccessKey && secretAccessKey.length > 1) {
    options.secretAccessKey = secretAccessKey;
}
const secretsmanager = new AWS.SecretsManager(options);

// Get Secret
const params = {
    SecretId: secret,
};

// Get secrets and create secret
const spinner = ora('Retrieving secret...').start();
secretsmanager.getSecretValue(params).promise()
    .then(({ SecretString }) => {
        spinner.stop();
        // Log
        console.log('Secrets retrieved.');
        const secrets = JSON.parse(SecretString);
        fs.writeFileSync(ymlFile, `---
apiVersion: v1
kind: Secret
metadata:
  name: "${secret.replace('/', '-')}"
type: Opaque
data:
${Object.keys(secrets).map((s) => {
            return `  ${s}: ${Buffer.from(secrets[s]).toString('base64')}`
        }).join("\n")}
---`, () => {
            console.log('Secret completed.')
        });
    })
    .catch((err) => {
        spinner.stop();
        console.error(`ERROR: Error: ${err.message}`);
        process.exit(1);
    });