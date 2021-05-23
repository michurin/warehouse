$(() => {
  const send = () => {
    const e = $('#text');
    const text = e.val();
    $.post('/chat/send', text);
    e.val(''); // TODO move it to `success` callback
    e.focus();
  };
  $('#button').click(send);
  $('#text').keyup(function (e) {
    if (e.which == 13) {
      send();
      return false;
    }
  });
  var lastID = 0;
  const poll = (id) => {
    if (id === undefined) {
      id = lastID;
    }
    $.post('/chat/poll', {id: id}, (data) => {
      const obj = JSON.parse(data);
      const id = obj.lastID;
      lastID = id;
      poll(id);
      if (obj.messages && obj.messages.forEach) {
        obj.messages.reverse().forEach((e) => { // we reverse original array here
          $('#conversation').append($('<div>').text(e.text));
        });
        const c = $('#conversation').children();
        for (var i = 0; i < c.length-10; i++) {
          c[i].remove();
        }
      }
    }).fail(() => {
      setTimeout(poll, 1000);
    });
  };
  poll();
  $('#text').focus();
});
