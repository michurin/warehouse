$(function() {
  var sender = chat({
    onmessages: function(messages) {
      messages.reverse().forEach(function(e) { // we reverse original array here
        $('#conversation').append($('<div>').text(e.text));
      });
      var c = $('#conversation').children();
      for (var i = 0; i < c.length-10; i++) {
        c[i].remove();
      }
    },
    onsuccess: function() {
      $('#text').val('').focus();
    },
  }),
   send = function() {
    sender($('#text').val()); // TODO object here
  };
  $('#button').click(send);
  $('#text').keyup(function(e) {
    if (e.which == 13) {
      send();
      return false;
    }
  });
  $('#text').focus();
});
