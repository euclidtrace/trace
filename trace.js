Office.onReady((info) => {
    if (info.host === Office.HostType.Excel) {
        document.getElementById('log-btn').addEventListener('click', handleButtonClick);
        document.getElementById('show-content-btn').addEventListener('click', showCellContent);
        document.getElementById('load-json-btn').addEventListener('click', loadJsonContent);
    }
});

function handleButtonClick() {
    console.log('Button clicked!');
}

async function loadJsonContent() {
    try {
        const response = await fetch('multidim_dag_resolution/gross_profit_calc_graph.json');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        document.getElementById('json-display').textContent = JSON.stringify(data, null, 2);
    } catch (error) {
        console.error('Error loading JSON:', error);
        document.getElementById('json-display').textContent = 'Error loading JSON: ' + error.message;
    }
}

async function showCellContent() {
    try {
        await Excel.run(async (context) => {
            const range = context.workbook.getSelectedRange();
            range.load("values");
            await context.sync();
            
            const display = document.getElementById("display");
            if (range.values && range.values.length > 0) {
                 display.innerText = range.values[0][0];
            } else {
                 display.innerText = "No content";
            }
        });
    } catch (error) {
        console.error(error);
        document.getElementById("display").innerText = "Error: " + error.message;
    }
}
