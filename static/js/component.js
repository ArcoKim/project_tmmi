if(sessionStorage.getItem("id") === null) {
    sessionStorage.setItem("id", crypto.randomUUID());
}

class Aside extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
        <aside>
            <div id="full-screen">
                <a href="/index.html">
                    <img src="/image/logo.png" alt="TMMI" id="logo">
                </a>
                <div id="mid-container">
                    <div class="side-container">
                        <a href="/index.html">
                            <div id="home">
                                <i class="fa-solid fa-house"></i>
                                <span>대시보드</span>
                                <hr>
                            </div>
                        </a>
                    </div>
                    <div class="side-container">
                        <h4>음악 검색</h4>
                        <ul class="icons">
                            <li>
                                <a href="/search/basic.html">
                                    <i class="fa-solid fa-magnifying-glass"></i>
                                    <span>일반 검색</span>
                                </a>
                            </li>
                            <li>
                                <a href="/search/advanced.html">
                                    <i class="fa-brands fa-searchengin"></i>
                                    <span>상세 검색</span>
                                </a>
                            </li>
                            <li>
                                <a href="/search/recommend.html">
                                    <i class="fa-solid fa-microscope"></i>
                                    <span>노래 추천</span>
                                </a> 
                            </li>
                        </ul>
                        <hr>
                    </div>
                    <div class="side-container">
                        <h4>차트 목록</h4>
                        <ul class="icons">
                            <li>
                                <a href="/chart/melon.html">
                                    <img src="/image/melon.jpeg" alt="Melon">
                                    <span>Melon</span>
                                </a>
                            </li>
                            <li>
                                <a href="/chart/genie.html">
                                    <img src="/image/genie.jpeg" alt="Genie">
                                    <span>Genie</span>
                                </a>
                            </li>
                            <li>
                                <a href="/chart/flo.html">
                                    <img src="/image/flo.jpeg" alt="FLO">
                                    <span>FLO</span>
                                </a>
                            </li>
                            <li>
                                <a href="/chart/bugs.html">
                                    <img src="/image/bugs.jpeg" alt="Bugs">
                                    <span>Bugs</span>
                                </a>
                            </li>
                            <li>
                                <a href="/chart/vibe.html">
                                    <img src="/image/vibe.jpeg" alt="Vibe">
                                    <span>Vibe</span>
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </aside>
        `;
    }
}

class Title extends HTMLElement {
    connectedCallback() {
        const h3 = document.createElement("h3");
        const main_text = document.createTextNode(this.getAttribute("main"));
        h3.appendChild(main_text);

        const p = document.createElement("p");
        const sub_text = document.createTextNode(this.getAttribute("sub"));
        p.appendChild(sub_text);

        this.appendChild(h3);
        this.appendChild(p);
    }
}

class Dashboard extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
        <div id="first-ranks">
            <div id="one-line">
                <div class="first" id="melon">
                    <div class="song-name">
                        <h4>Example</h4>
                        <p><span class="melon">Melon</span> 차트 1위 곡</p>
                    </div>
                    <div class="song-team">
                        <img src="https://cdn.icon-icons.com/icons2/510/PNG/512/image_icon-icons.com_50366.png">
                        <p>Arco</p>
                    </div>
                </div>
                <div class="first" id="genie">
                    <div class="song-name">
                        <h4>Example</h4>
                        <p><span class="genie">Genie</span> 차트 1위 곡</p>
                    </div>
                    <div class="song-team">
                        <img src="https://cdn.icon-icons.com/icons2/510/PNG/512/image_icon-icons.com_50366.png">
                        <p>Arco</p>
                    </div>
                </div>
                <div class="first" id="flo">
                    <div class="song-name">
                        <h4>Example</h4>
                        <p><span class="flo">FLO</span> 차트 1위 곡</p>
                    </div>
                    <div class="song-team">
                        <img src="https://cdn.icon-icons.com/icons2/510/PNG/512/image_icon-icons.com_50366.png">
                        <p>Arco</p>
                    </div>
                </div>
            </div>
            <div id="second-line">
                <div class="first" id="bugs">
                    <div class="song-name">
                        <h4>Example</h4>
                        <p><span class="bugs">Bugs</span> 차트 1위 곡</p>
                    </div>
                    <div class="song-team">
                        <img src="https://cdn.icon-icons.com/icons2/510/PNG/512/image_icon-icons.com_50366.png">
                        <p>Arco</p>
                    </div>
                </div>
                <div class="first" id="vibe">
                    <div class="song-name">
                        <h4>Example</h4>
                        <p><span class="vibe">Vibe</span> 차트 1위 곡</p>
                    </div>
                    <div class="song-team">
                        <img src="https://cdn.icon-icons.com/icons2/510/PNG/512/image_icon-icons.com_50366.png">
                        <p>Arco</p>
                    </div>
                </div>
            </div>
        </div>
        <div id="chart-list">
            <div class="chart">
                <canvas id="album"></canvas>
            </div>
            <div class="chart">
                <canvas id="artist"></canvas>
            </div>
        </div>
        `;
        
        fetch("/api/dashboard")
            .then((response) => response.json())
            .then((data) => {
                const song = data.song;
                for(let platform in song) {
                    const info = song[platform];
                    $(`#${platform} > div.song-name > h4`).text(info.name);
                    $(`#${platform} > div.song-team > p`).text(info.artist);
                    $(`#${platform} > div.song-team > img`).attr("src", info.image);
                }
                const album = data.album;
                new Chart(document.getElementById("album"), {
                    type: 'bar',
                    data: {
                      labels: album.label,
                      datasets: [{
                        label: 'Chart Album Point',
                        data: album.data,
                        borderWidth: 1
                      }]
                    },
                    options: {
                      scales: {
                        y: {
                          beginAtZero: true
                        }
                      }
                    }
                });
                const artist = data.artist;
                new Chart(document.getElementById("artist"), {
                    type: 'bar',
                    data: {
                      labels: artist.label,
                      datasets: [{
                        label: 'Chart Artist Point',
                        data: artist.data,
                        borderWidth: 1,
                        borderColor: 'rgb(0, 255, 0)',
                        backgroundColor: 'rgb(144, 238, 144)',
                      }]
                    },
                    options: {
                      scales: {
                        y: {
                          beginAtZero: true
                        }
                      }
                    }
                });
            });
    }
}

const chart = async (link, body={}) => {
    const plinks = {
        "melon": "https://www.melon.com/song/detail.htm?songId=", 
        "genie": "https://www.genie.co.kr/detail/songInfo?xgnm=", 
        "flo": "https://www.music-flo.com/detail/track/", 
        "bugs": "https://music.bugs.co.kr/track/", 
        "vibe": "https://vibe.naver.com/track/"
    };

    let options;
    if (Object.keys(body).length === 0) {
        options = {
            method: "GET"
        }
    } else {
        options = {
            method: "POST",
            body: JSON.stringify(body)
        }
    }

    return fetch(link, options)
        .then((response) => response.json())
        .then((chart) => {
            let tag = "";
            for(let [index, music] of chart.entries()) {
                let platforms = "";
                for(let platform in plinks) {
                    if(platform in music) {
                        let links = plinks[platform] + music[platform];
                        if(platform == "flo") {
                            links += "/details";
                        }
                        platforms += `
                        <a href="${links}">
                            <img src="/image/${platform}.jpeg" alt="${platform}"">
                        </a>
                        `;
                    }
                }
                tag += `
                <tr>
                    <td>${index + 1}</td>
                    <td>${music.name}</td>
                    <td>${music.artist}</td>
                    <td class="album">
                        <img src="${music.image}" alt="${music.album}">
                        <span>${music.album}</span>
                    </td>
                    <td class="platforms">
                        ${platforms}
                    </td>
                </tr>
                `;
            }
            return tag;
        });
};

let table = null;
class Charts extends HTMLElement {
    connectedCallback() {
        chart(`/api/chart?type=${this.getAttribute("platform")}&date=2024-06-13`)
            .then((tbody) => {
                this.innerHTML = `
                <div id="title">
                    <input type='date' placeholder='날짜를 선택해주세요.' value="2024-06-13" max="2024-06-13" id="date" required />
                </div>
                <table id="chart">
                    <thead>
                        <tr>
                            <th>순위</th>
                            <th>곡 이름</th>
                            <th>아티스트</th>
                            <th>앨범</th>
                            <th>플랫폼</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${tbody}
                    </tbody>
                    <tfoot>
                        <tr>
                            <th>순위</th>
                            <th>곡 이름</th>
                            <th>아티스트</th>
                            <th>앨범</th>
                            <th>플랫폼</th>
                        </tr>
                    </tfoot>
                </table>
                `;
                
                table = $("#chart").DataTable({autoWidth: false});
            });
    }
}

class Search extends HTMLElement {
    connectedCallback() {
        const types = this.getAttribute("type");
        let desc = "";
        if(types == "basic") {
            desc = "찾고 싶은 노래, 아티스트, 앨범을 입력해주세요.";
        }
        if(types == "advanced") {
            desc = "노래에 대한 간단한 설명을 입력해주세요.";
        }
        this.innerHTML = `
        <div class="wrap">
            <div id="search">
                <input type="text" id="searchTerm" placeholder="${desc}">
                <button type="submit" class="searchButton" onclick="${types}_search();">
                    <i class="fa fa-search"></i>
                </button>
            </div>
        </div>
        `;
    }
}

const basic_search = () => {
    const search = document.getElementById("searchTerm").value;
    chart("/api/search/basic?input=" + search)
        .then((tbody) => {
            table.destroy();
            document.querySelector("#result > tbody").innerHTML = tbody;
            table = $("#result").DataTable({autoWidth: false});
        });
};

class Basic extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
        <table id="result">
            <thead>
                <tr>
                    <th>#</th>
                    <th>곡 이름</th>
                    <th>아티스트</th>
                    <th>앨범</th>
                    <th>플랫폼</th>
                </tr>
            </thead>
            <tbody>
            </tbody>
            <tfoot>
                <tr>
                    <th>#</th>
                    <th>곡 이름</th>
                    <th>아티스트</th>
                    <th>앨범</th>
                    <th>플랫폼</th>
                </tr>
            </tfoot>
        </table>
        `;
        table = $("#result").DataTable({autoWidth: false});
    }
}

const advanced_search = () => {
    document.getElementsByClassName("dt-empty")[0].textContent = "Loading ...";
    const prompt = document.getElementById("searchTerm").value;
    chart("/api/search/advanced", {prompt: prompt})
        .then((tbody) => {
            table.destroy();
            document.querySelector("#result > tbody").innerHTML = tbody;
            table = $("#result").DataTable({autoWidth: false});
        });
};

class Advanced extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
        <table id="result">
            <thead>
                <tr>
                    <th>#</th>
                    <th>곡 이름</th>
                    <th>아티스트</th>
                    <th>앨범</th>
                    <th>플랫폼</th>
                </tr>
            </thead>
            <tbody>
            </tbody>
            <tfoot>
                <tr>
                    <th>#</th>
                    <th>곡 이름</th>
                    <th>아티스트</th>
                    <th>앨범</th>
                    <th>플랫폼</th>
                </tr>
            </tfoot>
        </table>
        `;
        table = $("#result").DataTable({autoWidth: false});
    }
}

const chat_time = () => {
    const today = new Date();
    let hour = today.getHours();
    hour = hour >= 10 ? hour : '0' + hour;
    let minute = today.getMinutes();
    minute = minute >= 10 ? minute : '0' + minute;
    return hour + ':' + minute
}

const bot_chat = (message) => {
    return `
    <div class="message incoming">
        <img src="https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEjK_KUycTsd1UDiUSNvxjcrb2McjD5Ov_PdW2dQjQPbqHoTnviRaKjD1ZUPJ9u1Z9AWyWA5EIfY-YxEF-ePwynZjiSGrlO3weBVKeBu1XNs_H0JprOzQFPiqsHnX-YD5ffQ1fO9CZGfyLEt9MHIGcS1qBeWBHCsqkcWgp9Suj2uo3xoLd-4xQfNCcS55Ms/s320/20240201_093035_0000.png" alt="Avatar" class="message-avatar">
        <div class="message-content">${message}</div>
        <div class="message-time">${chat_time()}</div>
    </div>
    `;
}

const recommend = () => {
    const prompt = document.getElementById("messageInput").value
    const chatCont = document.getElementById("chatContainer")

    chatCont.innerHTML = `
    <div class="message outgoing">
        <img src="https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEhwj8E__pXohJ70JOQlvENXnPedGn7SDArF0nA84pOfDxqicmX8iwdelVCIAB023O-Fo7ieKgCdNvwo1BT7u8-Q25os-9jOTcMTTNanPwwXbq-cuZEmkP0tuzy8nb17o94doiIFGm2sVkPwf8Wcb6Y-Gl0RaGK2qvweuoiH8363o9_bxHaQ5FrJ1USMwWY/s320/IMG_20240201_092300.jpg" alt="Avatar" class="message-avatar">
        <div class="message-content">${prompt}</div>
        <div class="message-time">${chat_time()}</div>
    </div>
    ` + chatCont.innerHTML;
    document.getElementById("messageInput").value = "";

    chatCont.innerHTML = bot_chat("Loading ...") + chatCont.innerHTML;
    fetch('/api/search/lyrics', {
        method: "POST",
        body: JSON.stringify({
            prompt: prompt,
            session_id: sessionStorage.getItem("id")
        })
    })
    .then((response) => response.text())
    .then((result) => {
        document.getElementsByClassName("message-content")[0].textContent = result;
    })
}

class Bot extends HTMLElement {
    connectedCallback() {
        this.innerHTML = `
        <div class="message-input">
            <input type="text" id="messageInput" placeholder="티미한테 전달할 메시지">
            <button onclick="recommend();"><i class="fas fa-paper-plane"></i></button>
        </div>
        <div class="chat-container" id="chatContainer">
            ${bot_chat("안녕하세요, 티미에요! 저한테 노래를 추천받으세요.")}
        </div>
        `;
    }
}

class Footer extends HTMLElement {
    connectedCallback() {
        const footer = document.createElement("footer");
        const p = document.createElement("p");
        const text = document.createTextNode("Copyright 2024. Kimjeongtae All rights reserved.");

        p.appendChild(text);
        footer.appendChild(p);

        this.appendChild(footer);
    }
}

customElements.define('tmmi-aside', Aside);
customElements.define('tmmi-title', Title);
customElements.define('tmmi-dashboard', Dashboard);
customElements.define('tmmi-chart', Charts);
customElements.define('tmmi-search', Search);
customElements.define('tmmi-basic', Basic);
customElements.define('tmmi-advanced', Advanced);
customElements.define('tmmi-bot', Bot);
customElements.define('tmmi-footer', Footer);