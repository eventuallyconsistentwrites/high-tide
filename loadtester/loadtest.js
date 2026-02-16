import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
    // We use the internal docker network alias 'high-tide-server'
    const url = `http://${__ENV.SERVER_HOST}:${__ENV.SERVER_PORT}/posts`;

    // Send request
    http.get(url);

    // Sleep for 2 seconds (same as your Go code)
    sleep(2);
}