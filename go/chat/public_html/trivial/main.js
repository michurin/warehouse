$(function () {
  var sender = window.chat({
    'onmessages': function (messages) {
      messages.reverse().forEach(function (msg) { // We reverse original array here
        $('#conversation').append($('<div>').text(msg));
      });
      var c = $('#conversation').children();
      for (var i = 0; i < c.length - 10; i++) {
        c[i].remove();
      }
    },
    'onsuccess': function () {
      $('#text').val('').focus();
    }
  });
  var send = function () {
    sender($('#text').val()); // String, however it can be object here
  };
  $('#button').click(send);
  $('#text').keyup(function (e) {
    if (e.which === 13) {
      send();
    }
  });
  $('#text').focus();
});
