import "./style.css";
import DataTable from "datatables.net-dt";
import "datatables.net-dt/css/dataTables.dataTables.css";

document.querySelector("#app").innerHTML = `
<div class="container">
    <div id="tracked" class="box">
        <h3>Total devices tracked</h3>
        <h1>-</h1>
    </div>
    <div id="unhealthy" class="box">
        <h3>Unhealthy devices</h3>
        <h1>-</h1>
    </div>
    <div id="offline" class="box">
        <h3>Offline devices</h3>
        <h1>-</h1>
    </div>
</div>
<div class="container">
    <table id="devices" style="width:1200px;margin: 20px auto;" class="cell-border stripe hover">
        <thead>
            <tr>
                <th>Site name</th>
                <th>Site description</th>
                <th>Device name</th>
                <th>Device description</th>
                <th>Problems</th>
            </tr>
        </thead>
        <tbody>
        </tbody>
    </table>
</div>
`;

let table = new DataTable("#devices", {
  paging: false,
});

fetchData();
setInterval(fetchData, 5000);

function printString(elementSelector, string) {
  const element = document.querySelector(elementSelector);
  if (element) {
    element.innerHTML = string;
  } else {
    console.error(`Element ${elementSelector} not found.`);
  }
}

async function fetchData() {
  try {
    const response = await fetch("/api/counts");
    const data = await response.text();
    const counts = data ? JSON.parse(data) : null;

    const unhealthy =
      counts.StorageDisrupted +
      counts.CpuOverutilized +
      counts.RamOverutilized +
      counts.StorageFull +
      counts.NetworkPacketLoss +
      counts.ImageHealthImpaired;

    printString("#tracked h1", counts.Tracked);
    printString("#unhealthy h1", unhealthy);
    printString("#offline h1", counts.Disconnected);
  } catch (error) {
    console.error("Error fetching counts:", error);
    printString("#tracked h1", "Error");
    printString("#unhealthy h1", "Error");
    printString("#offline h1", "Error");
  }

  try {
    const response = await fetch("/api/devices");
    const data = await response.text();
    const devices = data ? JSON.parse(data) : null;

    table.clear();

    devices.forEach((device) => {
      const problems = [];

      if (device.StorageDisrupted) problems.push("Storage disrupted");
      if (device.CpuOverutilized) problems.push("CPU overutilized");
      if (device.RamOverutilized) problems.push("RAM overutilized");
      if (device.StorageFull) problems.push("Storage full");
      if (device.NetworkPacketLoss) problems.push("Network packet loss");
      if (device.ImageHealthImpaired) problems.push("Image health impaired");
      if (device.Disconnected) problems.push("Disconnected");

      const url = "https://platform.yoursix.com/devices/edit/" + device.Id;

      table.row
        .add([
          device.SiteName,
          device.SiteDescription,
          `<a href="${url}">${device.Name}</a>`,
          device.Description,
          problems.join(", "),
        ])
        .draw();
    });
  } catch (error) {
    console.error("Error fetching devices:", error);
    table.clear().draw();
  }
}
