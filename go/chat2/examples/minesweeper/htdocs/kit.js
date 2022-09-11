const chatAdapter = () => {
  const call = (url, data) => {
    if (!data) {
      return undefined;
    }
    return fetch(url, {
      method: 'POST',
      cache: 'no-cache',
      body: JSON.stringify(data),
    });
  };
  const sleep = (delay) => new Promise((resolve) => { setTimeout(resolve, delay); });
  const adapter = {};
  adapter.loop = async (f) => {
    let bound = [];
    while (true) {
      try {
        const response = await call('/sub', { b: bound });
        if (response.status !== 200) {
          throw response.status;
        }
        const data = await response.json();
        bound = data.b || bound; // fallback to previous value if no updates
        f(data); // TODO multiply calls
      } catch (e) {
        console.log(e);
        await sleep(2000);
      }
    }
  };
  adapter.send = (m) => { call('/pub_chat', m); };
  adapter.game = (m) => { call('/pub_game', m); };
  return adapter;
};
