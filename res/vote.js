
var myData = {vote: "steve", other_vote: "4"};


myVotes = {};

function countVote(candidate) {
// All this converting stuff to strings is for the json unmarshaller
// Yeah, I'm not sure it's the best solution either
	candidate = candidate.toString()
	if( candidate in myVotes ){
		myVotes[candidate] = (1 + parseInt(myVotes[candidate])).toString();
	} else {
		myVotes[candidate] = "1";
	}
}


function SendData(sendMe) {
  $(function() {
  	// If there are no votes, dont bother sending data
	if( Object.keys(sendMe).length === 0 ) {
      $("#the_span").text("No clicks, no data sent");
	  return;
	}

    $.ajax({
      url: '/vote',
      type: 'post',
      dataType: 'json',
      contentType: 'application/json',
      data: JSON.stringify(sendMe),
      success: function( data, status, xhr) {
      	$("#the_span").text("Data: " + data['reponse']);
      },
      error: function(xhr, status, error) {
        $("#the_span").text("Status: " + status +", Error: " + error);
      }
    })
  });
}

countVote("5")

//  send data every 10 seconds
var send_interval = 10*1000;

function progress_bar() {
  var width = 1;
  var id = setInterval(frame, send_interval/100);
  function frame() {
    if(width >= 100) {
      clearInterval(id);
      SendData(myVotes);
      myVotes = {};
    } else {
      width++;
      $("#progressBar").css('width',width+'%');
    }
  }
}

progress_bar()
setInterval( progress_bar, send_interval);
