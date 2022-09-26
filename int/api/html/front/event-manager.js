//TODO 	Events must be managed by Operation Id
// The retry mechanism must be done by mimic the same operation of the operation ID related to the event in
class EventManager {
    async subscribe(str, address, callback) {
        axios({
            url: `/thyra/events/${str}/${address}`,
            method: "GET",
        })
            .then((resp) => {
                callback(resp);
            })
            .catch((e) => {
                // TODO Implement retry mechanism here
                console.error(e);
            });
    }
}
