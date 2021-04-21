import { getActiveEnv } from './env.js';

const socketUrl = 'ws://' + document.location.host + '/ws';
let conn;

const loadSocket = ({
    open = null,
    message = null,
    error = null,
    close = null,
}) => {
    conn = new WebSocket(socketUrl);
    if (open) {
        conn.onopen = open;
    }
    if (message) {
        conn.onmessage = message;
    }
    if (error) {
        conn.onerror = error;
    }
    if (close) {
        conn.onclose = close;
    }
};

const subscribe = (topicKeys) => {
    const { name } = getActiveEnv();
    const msg = {
        type: 'consume',
        env: name,
        payload: topicKeys,
    };
    conn.send(JSON.stringify(msg));
};

const unsubscribe = (topics) => {
    const { name } = getActiveEnv();
    const msg = {
        type: 'unsubscribe',
        env: name,

        payload: topics,
    };
    conn.send(JSON.stringify(msg));
};

export {
    loadSocket,
    subscribe,
    unsubscribe,
};