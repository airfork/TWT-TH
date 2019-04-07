'use strict';

// Initialize modal element
var modal = null;
document.addEventListener('DOMContentLoaded', function () {
    const modalElement = document.querySelectorAll('.modal');
    const instances = M.Modal.init(modalElement);
    modal = instances[0];
});
// Setup AJAX request to get JSON data
let request = new XMLHttpRequest();
request.open('GET', '/data', true);

request.onload = function () {
    if (request.status >= 200 && request.status < 400) {
        // Success!
        // Parse data
        const data = JSON.parse(request.responseText);
        const table = document.getElementById('table');
        let rows = '';
        // Iterate over data returned and created rows
        // that will eventually be added to the html
        for (let i = 0; i < data.length; i++) {
            const country = data[i];
            rows += `
            <tr class="country-row" style="cursor: pointer">
                <td>${country.Country}</td>
                <td>$${country.TotalSalesString}</td>
                <td>${country.BestSelling}</td>
                <td>${country.QuantitySoldString}</td>
            </tr> 
            `
        }
        // Add rows to page
        table.innerHTML = rows;
        const trs = document.querySelectorAll('.country-row');
        // Onclick, each row opens a modal with a breakdown of car manufacturers
        for (let i = 0; i < trs.length; i++) {
            trs[i].onclick = () => {
                showgraph(data[i]);
            }
        }
    } else {
        // We reached our target server, but it returned an error
        M.toast({
            html: 'There was an error of some sort, please try again.'
        });
        console.log('Error');
    }
};

request.onerror = function () {
    // There was a connection error of some sort
    M.toast({
        html: 'There was an error of some sort, please try again.'
    });
    console.log("There was an error of some sort, please try again");
};

request.send();

// Show graph takes in country data and renders a pie chart with it
function showgraph(country) {
    let datapoints = [];
    // Iterate over car makes
    for (let key in country.Makes) {
        if (key.trim() === '') {
            continue;
        }
        let y = (country.Makes[key] / country.QuantitySold * 100).toFixed(2);
        // Data points in canvasJS need to be in this format
        datapoints.push({y: y, label: key});
    }
    // Build chart and display modal
    var chart = new CanvasJS.Chart("chartContainer", {
        title: {
            text: `${country.Country} Breakdown`
        },
        backgroundColor: '#fafafa',
        animationEnabled: true,
        data: [{
            type: "pie",
            startAngle: 240,
            yValueFormatString: "##0.00\"%\"",
            indexLabel: "{label} {y}",
            dataPoints: datapoints,
            horizontalAlign: "center",
        }]
    });
    chart.render();
    modal.open();
}