// https://developer.chrome.com/docs/extensions/reference/webRequest/
chrome.webRequest.onBeforeRequest.addListener((details) => {
  if (details.type === 'script' && details.url.substr(-3) === '.js') {
    return { redirectUrl: `${details.url}?noRT` };
  }
}, { urls: ['http://*/*.js'] }, ['blocking']);
