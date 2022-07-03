$(() => {
  const kit = chatAdapter('/pub', '/sub');
  const text = $('#text');
  text.focus();
  text.keypress((e) => {
    if (e.which === 13) {
      kit.send(text.val()); // in fact, you are able to send any object
      text.val('');
      text.focus();
      return false;
    }
    return true;
  });
  kit.loop((messages) => {
    messages.forEach((msg) => {
      console.log(msg);
      $('#board').prepend($('<div>').text(msg));
    });
  });
});
