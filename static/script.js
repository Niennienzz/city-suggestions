(function () {
    'use strict';

    let input = document.querySelector('input');
    let container = document.getElementsByClassName('result')[0];
    let currentFocusIndex;

    input.addEventListener('input', (e) => {
        const value = e.target.value;
        if (!value) {
            container.style.display = 'none';
        }
        if (result.length) {
            buildResult(result, value);
        } else {
            container.style.display = 'none';
        }
    });

    input.addEventListener('input', debounce((e) => {
        const value = e.target.value;
        if (value) {
            currentFocusIndex = -1;
            fetch(`http://localhost:8000/city/search?q=${value}`).then(res => res.json()).then(result => {
                if (result.length) {
                    buildResult(result, value);
                } else {
                    container.style.display = 'none';
                }
            });
        }
    }, 10));

    input.addEventListener('keydown', (e) => {
        if (container.style.display === 'none') {
            return;
        }
        switch (e.keyCode) {
            case 40:
                currentFocusIndex++;
                updateTargetNode();
                break;
            case 38:
                currentFocusIndex--;
                updateTargetNode();
                break;
            case 13:
                if (currentFocusIndex >= 0) {
                    input.value = container.childNodes[currentFocusIndex].innerText;
                    cleanResult();
                }
            default:
                break;
        }
    });

    function buildResult(result, value) {
        cleanResult();

        let content = document.createDocumentFragment();
        result.forEach(item => {
            let div = document.createElement('div');
            div.style.padding = '5px';
            div.innerHTML = `<strong>${value}</strong>${item}`;
            div.addEventListener('click', () => {
                input.value = item;
            });
            content.appendChild(div);
        });

        container.style.display = 'block';
        container.appendChild(content);
    }

    function updateTargetNode() {
        const children = container.childNodes;
        const childNodeCount = container.childElementCount;
        if (currentFocusIndex < 0) {
            currentFocusIndex = childNodeCount - 1;
        }
        if (currentFocusIndex >= childNodeCount) {
            currentFocusIndex = 0;
        }
        children.forEach((child, idx) => {
            child.style.backgroundColor = idx === currentFocusIndex ? 'lightblue' : 'white';
        });
    }

    function debounce(fn, delay) {
        let timer;

        return function () {
            let context = this;
            let args = arguments;

            clearTimeout(timer);

            timer = setTimeout(() => {
                fn.apply(context, args);
            }, delay);
        }
    }

    function cleanResult() {
        while (container.firstChild) {
            container.removeChild(container.firstChild);
        }
        container.style.display = 'none';
    }

    document.addEventListener('click', () => {
        cleanResult();
    });
})();
