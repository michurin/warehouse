window.Room = function () {
  var nop = function () {};
  // Defaults
  this.room = '__default__';
  this.urlPoll = '/api/poll';
  this.urlPublish = '/api/publish';
  this.onsuccess = nop;
  this.onerror = nop;
  this.onmessages = nop;
  this.onconnectionup = nop;
  this.onconnectiondown = nop;
  // Closure
  var that = this;
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
  var isDown = false;
  var poll = function () {
    request(
      that.urlPoll,
      JSON.stringify({'id': lastID, 'room': that.room}),
      function (st, body) {
        if (st === 200) {
          window.setTimeout(poll, 0);
          var obj = JSON.parse(body);
          lastID = obj.lastID || 0;
          var msgs = obj.messages || [];
          if (msgs) {
            that.onmessages(msgs);
          }
          if (isDown) {
            isDown = false;
            that.onconnectionup();
          }
        } else {
          isDown = true;
          window.setTimeout(poll, 1000);
          that.onconnectiondown();
        }
      }
    );
  };
  this.send = function (msg) {
    request(
      that.urlPublish,
      JSON.stringify({'message': msg, 'room': that.room}),
      function (st, body) {
        if (st === 200) {
          that.onsuccess(body);
        } else {
          that.onerror();
        }
      }
    );
  };
  this.run = function () {
    poll(); // TODO? Check if already run
  };
};
