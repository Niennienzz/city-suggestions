document
    .querySelector("#autoComplete")
    .addEventListener("autoComplete", event => {
        console.log(event);
    });

const autoCompleteJS = new autoComplete({
    data: {
        src: async () => {
            const query = document.querySelector("#autoComplete").value;
            const source = await fetch(`http://localhost:8000/city/search?q=${query}`);
            return await source.json();
        },
        key: ["name"],
        cache: false,
    },
    placeHolder: "Search City",
    selector: "#autoComplete",
    threshold: 0,
    debounce: 25,
    searchEngine: "loose",
    highlight: true,
    maxResults: 25,
    resultsList: {
        render: true,
        container: source => {
            source.setAttribute("id", "autoComplete_list");
        },
        destination: document.querySelector("#autoComplete"),
        position: "afterend",
        element: "ul"
    },
    noResults: () => {
        const result = document.createElement("li");
        result.setAttribute("class", "no_result");
        result.setAttribute("tabindex", "1");
        result.innerHTML = "No Results";
        document.querySelector("#autoComplete_list").appendChild(result);
    },
    onSelection: feedback => {
        const selection = feedback.selection.value.name;
        document.querySelector(".selection").innerHTML = selection;
        document.querySelector("#autoComplete").value = "";
        document
            .querySelector("#autoComplete")
            .setAttribute("placeholder", selection);
        console.log(feedback);
    }
});
