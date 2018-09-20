module.exports = {
  echo: function (event, context) {
    console.log(`THE DATA: ${event.data}`);
    return event.data;
  }
}