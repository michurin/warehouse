window.chat = function (options) {
  var nop = function () {};
  // Options
  var onmessages = options.onmessages || nop;
  var ondown = options.ondown || nop;
  var onsuccess = options.onsuccess || nop;
  var onerror = options.onerror || nop;
  var publishUrl = options.publishUrl || '/api/publish';
  var pollUrl = options.pollUrl || '/api/poll';
  // TODO argument? option?
  var room = options.room || '__default__';
  // XHR inspired by jQuery
  var request = function (url, data, complete) {
    var xhr = new window.XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.timeout = 300000; // Default for most browsers
    xhr.onload = function () {
      complete(xhr.status, xhr.responseText);
      complete = nop;
    };
    xhr.onabort = xhr.onerror = xhr.ontimeout = function () {
      complete(xhr.status, undefined);
      complete = nop;
    };
    try {
      xhr.send(data);
    } catch (exc) {
      complete(599, undefined);
      complete = nop;
    }
  };
  // Polling
  var lastID = 0;
  var poll = function () {
    request(pollUrl, JSON.stringify({'id': lastID, 'room': room}), function (st, body) {
      if (st === 200) {
        var obj = JSON.parse(body);
        lastID = obj.lastID;
        window.setTimeout(poll, 0);
        onmessages(obj.messages || []);
      } else {
        ondown();
        window.setTimeout(poll, 1000);
      }
    });
  };
  poll();
  // Sender
  return function (message) {
    request(publishUrl, JSON.stringify({'message': message, 'room': room}), function (st, body) {
      if (st === 200) {
        onsuccess(body); // TODO id (has to be produced by server)
      } else {
        onerror();
      }
    });
  };
};
