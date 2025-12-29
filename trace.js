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
        const container = document.getElementById('json-display');
        container.innerHTML = ''; // Clear previous content
        container.className = 'box tree-view'; // Add tree-view class
        container.appendChild(createTreeView(data));
    } catch (error) {
        console.error('Error loading JSON:', error);
        document.getElementById('json-display').textContent = 'Error loading JSON: ' + error.message;
    }
}

function createTreeView(data, key = null) {
    // Handle primitive types (leaf nodes)
    if (data === null || typeof data !== 'object') {
        const span = document.createElement('span');
        if (key !== null) {
            const keySpan = document.createElement('span');
            keySpan.className = 'tree-key';
            keySpan.textContent = `"${key}": `;
            span.appendChild(keySpan);
        }
        const valSpan = document.createElement('span');
        valSpan.className = `tree-${data === null ? 'null' : typeof data}`;
        valSpan.textContent = JSON.stringify(data);
        span.appendChild(valSpan);
        return span;
    }

    // Handle Objects and Arrays (branch nodes)
    const container = document.createElement('div');
    
    const header = document.createElement('div');
    const toggle = document.createElement('span');
    toggle.className = 'tree-toggle';
    toggle.textContent = '[-]'; // Default expanded
    
    const label = document.createElement('span');
    if (key !== null) {
        const keySpan = document.createElement('span');
        keySpan.className = 'tree-key';
        keySpan.textContent = `"${key}": `;
        label.appendChild(keySpan);
    }
    
    const isArray = Array.isArray(data);
    const openBrace = document.createElement('span');
    openBrace.textContent = isArray ? '[' : '{';
    
    header.appendChild(toggle);
    header.appendChild(label);
    header.appendChild(openBrace);
    container.appendChild(header);

    const children = document.createElement('ul');
    const keys = Object.keys(data);
    
    keys.forEach((k, i) => {
        const li = document.createElement('li');
        // Pass key only if it's an object, for arrays we just show values usually, 
        // but to match JSON structure strictly we can omit key for array items or show index.
        // Standard JSON view: Object keys are shown, Array indices are usually implied.
        const childNode = createTreeView(data[k], isArray ? null : k);
        li.appendChild(childNode);
        if (i < keys.length - 1) {
            li.appendChild(document.createTextNode(','));
        }
        children.appendChild(li);
    });

    const closeBrace = document.createElement('div');
    closeBrace.textContent = isArray ? ']' : '}';
    closeBrace.style.paddingLeft = '24px'; // Indent closing brace

    container.appendChild(children);
    container.appendChild(closeBrace);

    // Toggle logic
    toggle.onclick = (e) => {
        e.stopPropagation();
        const isHidden = children.style.display === 'none';
        if (isHidden) {
            children.style.display = 'block';
            closeBrace.style.display = 'block';
            toggle.textContent = '[-]';
            openBrace.textContent = isArray ? '[' : '{';
        } else {
            children.style.display = 'none';
            closeBrace.style.display = 'none';
            toggle.textContent = '[+]';
            // Show summary when collapsed
            const count = keys.length;
            openBrace.textContent = isArray ? `[ ${count} items ]` : `{ ... }`;
        }
    };

    return container;
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
