$(() => {
  const kit = chatAdapter('/pub', '/sub');
  const text = $('#text');
  const name = $('#name');
  const color = $('#color');
  text.keypress((e) => {
    if (e.which === 13) {
      kit.send({
        color: color.val(),
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
  text.focus();
  kit.loop((messages) => {
    messages.forEach((msg) => {
      $('#board').append($('<div>').append(
        $('<b>').text(`${msg.name}:`),
        $('<span>').text(` ${msg.text}`),
      ).css({ color: msg.color }));
    });
    $('html, body').scrollTop($(document).height() - $(window).height());
  });
});
