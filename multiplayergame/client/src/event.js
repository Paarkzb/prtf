/**
 * Event is used to wrap all messages Send and Receive
 * on the Websocket
 * The type is used as a RPC
 **/
export default class EventMsg {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}
