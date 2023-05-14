$(document).ready(function () {

    $.get("/api/random/joke", function (data) {
        $("#dadJoke").html(data.joke);
    });

});