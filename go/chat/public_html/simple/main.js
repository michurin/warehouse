$(function () {

  // Common

  var currentFigure = 'x';

  var setFigure = function (f) {
    // TODO produce log-message
    currentFigure = f;
    $('.fig').removeClass('sel');
    $('#fig-' + f).addClass('sel');
  };

  // Setup env

  var setupShareLink = function () {
    var f = function () {
      try {
        var s = window.getSelection();
        var r = window.document.createRange();
        r.selectNode(this);
        s.removeAllRanges();
        s.addRange(r);
      } catch (exc) {}
    };
    $('#share').text(window.location.href).hover(f);
  };

  var room = window.location.hash.substr(1);
  if (room === '') {
    room = Math.ceil(Math.random() * 46655).toString(36);
    window.location.hash = room;
  }
  setupShareLink();

  // Setup chat

  var chat = new Room();
  chat.room = room;
  chat.urlPoll = '/api/simple/poll';
  chat.urlPublish = '/api/simple/publish';

  // Setup and run receiver loop

  chat.onsuccess = function () { // TODO (e)
    $('#text').val('').focus();
  };
  chat.onconnectiondown = function () { $('#error').show(); };
  chat.onconnectionup = function () { $('#error').hide(); };

  chat.run();

  // Setup sender

  var cellId = function (x, y) {
    return 'E' + x + '-' + y;
  };

  var initGameArea = function (size) {
    var mkHandler = function (x, y) {
      return function (e) {
        e.preventDefault();
        chat.send({
          // TODO add nick, color and produce log-message
          'type': 'game-tune',
          'x': x,
          'y': y,
          'fig': currentFigure
        });
      };
    };
    var t = $('#game-area');
    t.empty();
    for (var i = 0; i < size; i++) {
      var tr = $('<tr>');
      for (var j = 0; j < size; j++) {
        var txt = '<img src="e.svg" id="' + cellId(j, i) + '">';
        var img = $(txt).click(mkHandler(j, i));
        var td = $('<td>').addClass('empty').append(img);
        tr.append(td);
      }
      t.append(tr);
    }
  };

  chat.onmessages = function (messages) {
    messages.reverse().forEach(function (msg) { // We reverse original array here
      if (msg.type === 'chat-message') {
        $('#conversation').append($('<div>').append(
          $('<span>').css('color', msg.color).text(msg.nick),
          $('<span>').html(':&nbsp;'),
          $('<span>').text(msg.text)
        ));
      } else if (msg.type === 'game-reset') {
        // TODO log message
        initGameArea(msg.size);
      } else if (msg.type === 'game-tune') {
        var e = $('#' + cellId(msg.x, msg.y));
        if (e.parent().hasClass('empty')) {
          // TODO log message
          e.parent().removeClass('empty');
          e.attr('src', msg.fig + '.svg');
        }
      }
    });
    var c = $('#conversation').children();
    for (var i = 0; i < c.length - 10; i++) {
      c[i].remove();
    }
  };

  // Setup UI

  var sendChatMessage = function () {
    chat.send({
      'type': 'chat-message',
      'text': $('#text').val(),
      'nick': $('#nick').val(),
      'color': $('#color').val()
    });
  };
  var sendGameReset = function (size) {
    chat.send({
      // TODO add nick, color and produce log-message
      'type': 'game-reset',
      'size': size
    });
  };

  $('#button').click(sendChatMessage);
  $('#text').keyup(function (e) {
    if (e.which === 13) {
      sendChatMessage();
    }
  });
  $('#text').focus();

  $('#start3').click(function () { sendGameReset(3); });
  $('#start15').click(function () { sendGameReset(15); });
  $('#fig-x').click(function () { setFigure('x'); });
  $('#fig-o').click(function () { setFigure('o'); });
  initGameArea(15);
  setFigure(Math.random() > 0.5 ? 'x' : 'o');
});
