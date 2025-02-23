'use strict';

chrome.declarativeNetRequest.onRuleMatchedDebug.addListener((e) => {
  const r = e.request;
  console.log(`Navigation blocked to ${r.url} method ${r.method} on tab ${r.tabId}. Type: ${r.type}. Initiator: ${r.initiator}.`);
});

console.log('Worker started.');
