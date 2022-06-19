const chatAdapter = (pubUrl, subUrl) => {
  const call = (url, data) => {
    return fetch(url, {
      method: 'POST',
      cache: 'no-cache',
      body: JSON.stringify(data),
    });
  };
  const sleep = delay => new Promise(r => setTimeout(r, delay));
  return {
    send(message) {
      call(pubUrl, message);
    },
    async loop(f) {
      var bound = -1;
      for (; ;) {
        try {
          const response = await call(subUrl, { bound });
          if (response.status != 200) {
            throw response.status;
          }
          const data = await response.json();
          bound = data.bound;
          for (i = data.messages.length - 1; i >= 0; i--) {
            f(data.messages[i]);
          }
        } catch {
          await sleep(2000);
        }
      }
    },
  };
};
