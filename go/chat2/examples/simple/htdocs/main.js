$(() => {
  const kit = chatAdapter('/pub', '/sub');
  const text = $('#text');
  const name = $('#name');
  text.keypress((e) => {
    if (e.which === 13) {
      kit.send({
        text: text.val(),
        name: name.val(),
      });
      text.val('');
      text.focus();
      return false;
    }
    return true;
  });
  name.keypress((e) => {
    if (e.which === 13) {
      text.focus();
      return false;
    }
    return true;
  });
  name.focus();
  kit.loop((messages) => {
    messages.reverse().forEach((msg) => {
      $('#board').append($('<div>').append(
        $('<b>').text(`${msg.name}:`),
        $('<span>').text(` ${msg.text}`),
      ));
    });
    $('html, body').scrollTop($(document).height() - $(window).height());
  });
});
