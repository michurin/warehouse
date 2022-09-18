const chatAdapter = () => {
  const sleep = (delay) => new Promise((resolve) => { setTimeout(resolve, delay); });
  const client = async (url, data) => {
    if (!data) { // skip empty requests, hackish
      return {};
    }
    try {
      const response = await fetch(url, {
        method: 'POST',
        cache: 'no-cache',
        body: JSON.stringify(data),
      });
      if (response.status !== 200) {
        throw response.status;
      }
      return await response.json();
    } catch (e) {
      console.log(e);
      await sleep(2000);
    }
    return {};
  }
  const adapter = {};
  adapter.loop = async (f) => {
    let bound = [];
    while (true) {
      try {
        const data = await client('/sub', { b: bound });
        bound = data.b || bound; // fallback to previous value if no updates
        f(data); // TODO multiply calls
      } catch (e) {
        console.log(e);
        await sleep(2000);
      }
    }
  };
  adapter.send = (m) => { return client('/pub_chat', m); };
  adapter.game = (m) => { return client('/pub_game', m); };
  return adapter;
};
