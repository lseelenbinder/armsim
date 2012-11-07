$(document).ready(function () {

  // Setup WebSockets
  ws = new WebSocket("ws://" + document.URL.split("/")[2] + "/ws");
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
    case "filename":
      $("#filename").text("File: " + received.Content);
      break;
    case "output":
      output(received);
      break;
    case "error":
      error(received);
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

  $("#terminal textarea").on("keydown", keyboardInput);

  // Setup Shortcut Keys
  $.Shortcuts.add({
    type: "down",
    mask: "Ctrl+O, O, Shift+O",
    handler: loadFile
  });
  $.Shortcuts.add({
    type: "down",
    mask: "Shift+?",
    handler: function() {
      $('#shortcuts').modal('show')
    }
  });
  $.Shortcuts.add({
    type: "down",
    mask: "F5",
    handler: function() {
      ws.send("start");
    }
  });
  $.Shortcuts.add({
    type: "down",
    mask: "F10",
    handler: function() {
      ws.send("step");
    }
  });
  $.Shortcuts.add({
    type: "down",
    mask: "Ctrl+B, B, Shift+B",
    handler: function() {
      ws.send("stop");
    }
  });
  $.Shortcuts.add({
    type: "down",
    mask: "Ctrl+T, T, Shift+T",
    handler: toggleTrace
  });
  $.Shortcuts.add({
    type: "down",
    mask: "Ctrl+R, R, Shift+R",
    handler: function() {
      ws.send("reset");
    }
  });
  $.Shortcuts.start();

  // Load Button
  $("#load-button").click(loadFile);

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

  $("#trace-button").click(toggleTrace);
  $("#system-trace-button").click(toggleSystemTrace);

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
  updateDisassembly(data.Disassembly, data.Registers[15]);
  updateRegisters(data.Registers);
  updateStack(data.Stack, data.Registers[13]);
  updateMemory(data.Memory);
  updateChecksum(data.Checksum);
  updateMode(data.Mode);
}

function updateChecksum(checksum) {
  $("#checksum").text("Checksum: " + checksum);
}

function updateMode(mode) {
  $("#mode").text("Mode: " + mode);
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

function updateStack(stack, sp) {
  $("#stack tbody").empty();
  $.each(stack, function (i) {
    $("#stack tbody").append("<tr><td>" + hexToString(sp) + "</td><td>" + hexToString(stack[i]) + "</td></tr");
    sp += 0x4;
  });
}

function updateDisassembly(instructions, pc) {
  $("#instructions").empty();
  var address = pc - 8;
  $.each(instructions, function (i) {
    if (instructions[i] == "") {
      return;
    }
    var encoded = parseInt(instructions[i].split("||")[0], 16);
    var decoded = instructions[i].split("||")[1].split(" ")[0] + "\t";
    var arguments = instructions[i].split("||")[1].split(" ").slice(1).join(" ");
    var comments = "";
    if (address == pc) {
      var active = "alert alert-success"
    } else {
      var active = "";
    }
    $("#instructions").append(
      "<div class='instruction " + active + "'><span class='address'>" + hexToString(address) +
      "</span><span class='encoded'>" + hexToString(encoded) +
      "</span><span class='decoded'>" + decoded +
      "<span class='arguments'>" + arguments +
      "</span></span><span class='comment'>" + comments + "</span></div>");
    address += 4
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

function toggleTrace() {
  var content = "on";
  if ($("#trace-button i").hasClass("icon-eye-close")) {
    // Turn tracing off
    content = "off";
  }

  ws.send("trace", content);

  $("#trace-button").html("<i class='icon-eye-"
    + ((content == "on") ? "close" : "open") +
    "'></i> Turn-" + (content == "on" ? "off" : "on") + " Tracing");
  $("#trace-button").toggleClass("btn-danger");
  $("#trace-button").toggleClass("btn-success");
}

function toggleSystemTrace() {
  var content = "off";
  if ($("#system-trace-button").hasClass("btn-success")) {
    // Turn system tracing off
    content = "on";
  }

  ws.send("system-trace", content);

  $("#system-trace-button").html("Turn-" + (content == "on" ? "off" : "on") + " System Tracing");
  $("#system-trace-button").toggleClass("btn-danger");
  $("#system-trace-button").toggleClass("btn-success");
}

function loadFile() {
  var filePath = prompt("Please enter your filename (relative to the executable).");

  disableButton("load");
  $(this).data("html", $(this).html());
  $(this).text("Loading...");

  $("#filename").text("File: " + filePath);

  ws.send("load", filePath);
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

function keyboardInput(e) {
  var c = String.fromCharCode(e.keyCode);
  ws.send("input", c);
  console.log("Sent: " + e.keyCode);
  return false;
}

function output(text) {
  if (text.Content[0] == '\r') {
    text.Content = "\n";
  }
  var old = $("#terminal textarea").val();
  $("#terminal textarea").val(old + text.Content);
}

function error(error) {
  alert("Error: " + error.Content);
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
