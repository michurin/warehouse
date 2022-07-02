const chatAdapter = (pubUrl, subUrl) => {
  const call = (url, data) => fetch(url, {
    method: 'POST',
    cache: 'no-cache',
    body: JSON.stringify(data),
  });
  const sleep = (delay) => new Promise((resolve) => { setTimeout(resolve, delay); });
  return {
    send(message) {
      call(pubUrl, message);
    },
    async loop(f) {
      let bound = 0;
      while (true) {
        try {
          const response = await call(subUrl, { bound });
          if (response.status !== 200) {
            throw response.status;
          }
          const data = await response.json();
          bound = data.bound;
          f(data.messages);
        } catch (e) {
          await sleep(2000);
        }
      }
    },
  };
};
