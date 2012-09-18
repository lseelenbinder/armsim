$(document).ready(function () {
  var ws = new WebSocket("ws://" + document.URL.split("/")[2] + "/ws");
  ws.onopen = function() {
    ws.send("hello");
  };

  ws.onmessage = function(e) {
    var received = JSON.parse(e.data);
    switch(received.Type) {
    case "status":
      switch(received.Content) {
      case "ready":
        ready();
        break;
      case "running":
        running();
        break;
      case "loaded":
        loaded();
        break;
      case "stopped":
        finished();
        break;
      case "finished":
        finished();
        break;
      }
      break;
    case "update":
      update(received);
      break;
    }
  };

  ws.onclose = function(e) {
    alert("connection closed")
  };

  ws._send = ws.send
  ws.send = function(type, data) {
    if (data === null) {
      ws._send($.toJSON({type: type, content: ""}));
    } else {
      ws._send($.toJSON({type: type, content: data}));
    }
  }


	$("#terminal .term").terminal(function(command, term) {

	}, {
		enabled: false,
		greetings: "ARMSim by Luke Seelenbinder"
	}
	);

  $("#load-button").click(function() {
    var filePath = prompt("Please enter your filename (relative to the executable).");

    disableButton("load");
    $(this).data("html", $(this).html());
    $(this).text("Loading...");

    $("#filename").text("File: " + filePath);

    ws.send("load", filePath);
  });

  $("#start-button").click(function() {
    ws.send("start");
  });

  $("#stop-button").click(function() {
    ws.send("stop");
  });

  $("#step-button").click(function() {
    ws.send("step");
  });

  $("#reset-button").click(function() {
    ws.send("reset");
  });

  $("#trace-button").click(function(e) {
    var content = "on";
    if ($("#trace-button i").hasClass("icon-eye-close")) {
      // Turn tracing off
      content = "off";
    }

    ws.send("trace", content);

    $("#trace-button").html("<i class='icon-eye-"
      + ((content == "on") ? "close" : "open") +
      "'></i> Turn-" + (content == "on" ? "off" : "on") + " Tracking");
    $("#trace-button").toggleClass("btn-danger");
    $("#trace-button").toggleClass("btn-success");
  });

  $("#memory-search").submit(function(e) {
    e.preventDefault();
    $(".memory-row").removeClass("active");
    var q = $(e.target).find("input[name='q']").val();
    if (q == "") {
      return
    }

    var address = parseInt(q, 16);
    var row = address >> 4;
    var kth = address & 0xF;

    row = $(".memory-row")[row];
    $("#memory-container").scrollTo( row );
    $(row).addClass("active");
  });
});

function update(data) {
  data = JSON.parse(data.Content);

  updateFlags(data.Flags);
  updateRegisters(data.Registers);
  updateMemory(data.Memory);
  updateChecksum(data.Checksum);
}

function updateChecksum(checksum) {
  $("#checksum").text("Checksum: " + checksum);
}

function updateFlags(flags) {
  $("#flags i").each(function (i, ele) {
    if (flags[i]) {
      $(ele).addClass("active");
      $(ele).removeClass("hidden");
    } else {
      $(ele).addClass("hidden");
      $(ele).removeClass("active");
    }
  });
}

function updateRegisters(registers) {
  $("#registers tbody").empty();
  $.each(registers, function (i) {
    $("#registers tbody").append("<tr><td>r" + i + "</td><td>" + hexToString(registers[i]) + "</td></tr");
  });
}

function updateMemory(memory) {
  $("#memory-container").empty();
  var row = "";
  var decoded = "";
  for (i = 0; i < memory.length; i++) {
    if (i % 16 == 0) {
      if (row.length > 0) {
        row += decoded + "</div>";
        $("#memory-container").append(row);
        decoded = "";
      }

      row = "<div class='memory-row'><span class='address'>" +
        hexToString(i) + "</span>";
    }
    row += nToWidth(memory[i].toString(16), 2) + " ";
    var ascii = String.fromCharCode(memory[i]);
    decoded += ((/[\x00-\x1F\x80-\xFF]/.test(ascii)) ? "." : ascii);
  }
}

function ready() {
  enableButton("load");
}

function running() {
  $.each(["start", "load", "step", "reset"], function(i, button) {
    disableButton(button);
  });
  enableButton("stop");
}

function loaded() {
  $.each(["start", "load", "step", "reset"], function(i, button) {
    enableButton(button);
  });

  $("#load-button").html($("#load-button").data("html"))
}

function finished() {
  $.each(["start", "load", "step", "reset"], function(i, button) {
    enableButton(button);
  });
  disableButton("stop");
}

function enableButton(button) {
  button = "#" + button + "-button"
  $(button).removeClass("disabled");
  $(button).removeAttr("disabled");
}

function disableButton(button) {
  button = "#" + button + "-button"
  $(button).addClass("disabled");
  $(button).attr("disabled", "disabled");
}

function hexToString(n) {
  n = n.toString(16);
  n = nToWidth(n, 8)

  return "0x" + n;
}

function nToWidth(n, width) {
  if (n.length < width) {
    var times = width - n.length
    for (j = 0; j < times; j++) {
      n = "0" + n;
    }
  }
  return n
}
