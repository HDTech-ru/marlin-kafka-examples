const {Kafka} = require('kafkajs');
const axios= require('axios');

const brokers = 'kafka1.example.org:443,kafka2.example.org:443'
const clientId = 'oauth-client-id';
const clientSecret = 'oauth-client-secret';
const tokenEndpoint = 'https://keycloak.example.org/auth/realms/kafka-authz/protocol/openid-connect/token'
const topic = 'example';

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

    const producer = kafka.producer();
    await producer.connect();

    const metadata = await producer.send({
        topic: topic,
        messages: [{key: "example_key", value: "Example message"}],
    });

    console.log(`Message was sent to the partition ${metadata[0].partition}. New offset is ${metadata[0].baseOffset}`);

    // Close the producer
    await producer.disconnect();
}

main().catch(console.error);
