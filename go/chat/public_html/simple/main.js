function cellId(x, y) {
  return 'E' + x + '-' + y;
}

var currentFigure = 'x';

function setFigure(f) {
  // TODO produce log-message
  currentFigure = f;
  $('.fig').removeClass('sel');
  $('#fig-' + f).addClass('sel');
}

$(function() {
  var sender = chat({
    publishUrl: '/api/small/publish',
    pollUrl: '/api/small/poll',
    onmessages: function(messages) {
      messages.reverse().forEach(function(msg) { // we reverse original array here
        if (msg.type == 'chat-message') {
          $('#conversation').append($('<div>').append(
            $('<span>').css('color', msg.color).text(msg.nick),
            $('<span>').html(':&nbsp;'),
            $('<span>').text(msg.text)
          ));
        } else if (msg.type == 'game-reset') {
          // TODO log message
          initGameArea(msg.size);
        } else if (msg.type == 'game-tune') {
          var e = $('#' + cellId(msg.x, msg.y));
          if (e.parent().hasClass('empty')) {
            // TODO log message
            e.parent().removeClass('empty');
            e.attr('src', msg.fig + '.svg');
          }
        }
      });
      var c = $('#conversation').children();
      for (var i = 0; i < c.length-10; i++) {
        c[i].remove();
      }
    },
    onsuccess: function() { // TODO (e)
      $('#text').val('').focus();
    },
  });

  var sendChatMessage = function() {
    sender({
      'type': 'chat-message',
      'text': $('#text').val(),
      'nick': $('#nick').val(),
      'color': $('#color').val()
    });
  };
  var sendGameReset = function(size) {
    sender({
      // TODO add nick, color and produce log-message
      'type': 'game-reset',
      'size': size,
    });
  }

  var initGameArea = function(size) {
    var t = $('#game-area');
    t.empty();
    var i, j;
    for (i = 0; i < size; i++) {
      var td = $('<tr>');
      for (j = 0; j < size; j++) {
        var id = 'E' + i + '-' + j;
        td.append($('<td>').addClass('empty').append($('<img src="e.svg" width="30" height="30" id="'+cellId(j, i)+'">').click(function(x, y){
          return function(e) {
            e.preventDefault();
            sender({
              // TODO add nick, color and produce log-message
              'type': 'game-tune',
              'x': x,
              'y': y,
              'fig': currentFigure,
            });
          };
        }(j, i))));
      }
      t.append(td);
    }
  }

  $('#button').click(sendChatMessage);
  $('#text').keyup(function(e) {
    if (e.which == 13) {
      sendChatMessage();
    }
  });
  $('#text').focus();

  $('#start3').click(function() {sendGameReset(3);});
  $('#start15').click(function() {sendGameReset(15);});
  $('#fig-x').click(function() {setFigure('x');});
  $('#fig-o').click(function() {setFigure('o');});
  initGameArea(15);
  setFigure(Math.random() > .5 ? 'x' : 'o');
});
