module.exports = {
  echo: function (event, context) {
    dataString = JSON.stringify(event.data)
    console.log(`THE DATA: ${dataString}`);
    return JSON.stringify(dataString);
  }
}