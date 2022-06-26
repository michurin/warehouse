$(() => {
  const kit = chatAdapter('/pub', '/sub');
  const text = $('#text');
  const name = $('#name');
  text.keypress((e) => {
    if (e.which == 13) {
      kit.send({
        text: text.val(),
        name: name.val(),
      });
      text.val('');
      text.focus();
      return false;
    }
  });
  name.keypress((e) => {
    if (e.which == 13) {
      text.focus();
      return false;
    }
  });
  name.focus();
  kit.loop((message) => {
    $('#board').append($('<div>').append(
      $('<b>').text(message.name + ':'),
      $('<span>').text(' ' + message.text)
    ));
    $("html, body").scrollTop($(document).height() - $(window).height());
  });
});
