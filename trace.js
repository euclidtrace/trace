Office.onReady((info) => {
    if (info.host === Office.HostType.Excel) {
        document.getElementById('log-btn').addEventListener('click', handleButtonClick);
        document.getElementById('show-content-btn').addEventListener('click', showCellContent);
    }
});

function handleButtonClick() {
    console.log('Button clicked!');
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
