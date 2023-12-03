const axios = require("axios");
const {Kafka} = require("kafkajs");

const brokers = 'kafka1.example.org:443,kafka2.example.org:443'
const clientId = 'oauth-client-id';
const clientSecret = 'oauth-client-secret';
const tokenEndpoint = 'https://keycloak.example.org/auth/realms/kafka-authz/protocol/openid-connect/token'
const topic = 'example';
const groupId = 'example';

async function getAccessToken() {
    return new Promise((resolve, reject) => {
        axios.post(tokenEndpoint, new URLSearchParams({
            grant_type: 'client_credentials',
            client_id: clientId,
            client_secret: clientSecret
        })).then(res => {
            resolve(res.data.access_token);
        }, err => {
            console.log('Error: ', err.message);
            reject(err);
        });
    });
}

function sleep(ms) {
    return new Promise((resolve) => {
        setTimeout(resolve, ms);
    });
}

async function main() {
    // Set up the Kafka client
    const kafka = new Kafka({
        clientId: 'marlin-kafka-example-app',
        brokers: brokers.split(','),
        ssl: true,
        sasl: {
            mechanism: 'oauthbearer',
            oauthBearerProvider: async () => {
                return {value: await getAccessToken()}
            }
        },
    });

    const consumer = kafka.consumer({groupId});

    await consumer.connect();

    await consumer.subscribe({ topic, fromBeginning: true })
    await consumer.run({
        eachMessage: async ({ partition, message }) => {
            console.log(`- [${partition} | ${message.offset}] ${message.key}#${message.value}`);
        },
    })

    await sleep(10000);

    // Close the producer
    await consumer.disconnect();
}

main().catch(console.error);
