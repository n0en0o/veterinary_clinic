const API_BASE = "http://localhost:8000/api";

document.addEventListener("DOMContentLoaded", function () {
  loadPets();
  loadOwners();
  loadOwnerSelect();
  loadPetSelect();

  document
    .getElementById("add-pet-form")
    .addEventListener("submit", handleAddPet);

  document
    .getElementById("add-record-form")
    .addEventListener("submit", handleAddHealthRecord);

  document
    .getElementById("add-owner-form")
    .addEventListener("submit", handleAddOwner);
});

//навигация
function showSection(sectionName) {
  document.querySelectorAll(".section").forEach((section) => {
    section.classList.remove("active");
  });

  switch (sectionName) {
    case "pets":
      document.getElementById("pets-section").classList.add("active");
      loadPets();
      break;
    case "owners":
      document.getElementById("owners-section").classList.add("active");
      loadOwners();
      break;
    case "add-owner":
      document.getElementById("add-owner-section").classList.add("active");
      break;
    case "add-pet":
      document.getElementById("add-pet-section").classList.add("active");
      break;
    case "add-record":
      document.getElementById("add-record-section").classList.add("active");
      break;
  }
}

//загрузка питомцев
async function loadPets() {
  try {
    const response = await fetch(`${API_BASE}/pets`);
    const pets = await response.json();
    displayPets(pets);
  } catch (error) {
    console.error("Error loading pets:", error);
  }
}

function displayPets(pets) {
  const container = document.getElementById("pets-list");
  container.innerHTML = pets
    .map(
      (pet) => `
        <div class="card pet-card" onclick="showPetDetail(${pet.id})">
            <h3>${pet.name}</h3>
            <p><strong>Вид:</strong> ${pet.species}</p>
            <p><strong>Порода:</strong> ${pet.breed || "Не указана"}</p>
            <p><strong>Владелец:</strong> ${pet.owner_name}</p>
            <p><strong>Возраст:</strong> ${calculateAge(pet.date_of_birth)}</p>
        </div>
    `
    )
    .join("");
}

// та же загрузка владельцев
async function loadOwners() {
  try {
    const response = await fetch(`${API_BASE}/owners`);
    const owners = await response.json();
    displayOwners(owners);
  } catch (error) {
    console.error("Error loading owners:", error);
  }
}

function displayOwners(owners) {
  const container = document.getElementById("owners-list");
  container.innerHTML = owners
    .map(
      (owner) => `
        <div class="card owner-card" onclick="showOwnerDetail(${owner.id})">
            <h3>${owner.first_name} ${owner.last_name}</h3>
            <p><strong>Email:</strong> ${owner.email}</p>
            <p><strong>Телефон:</strong> ${owner.phone || "Не указан"}</p>
            <p><strong>Адрес:</strong> ${owner.address || "Не указан"}</p>
        </div>
    `
    )
    .join("");
}

//загрузка выпадающих списков
async function loadOwnerSelect() {
  try {
    const response = await fetch(`${API_BASE}/owners`);
    const owners = await response.json();
    const select = document.getElementById("owner-select");
    select.innerHTML =
      '<option value="">Выберите владельца</option>' +
      owners
        .map(
          (owner) =>
            `<option value="${owner.id}">${owner.first_name} ${owner.last_name}</option>`
        )
        .join("");
  } catch (error) {
    console.error("Error loading owners for select:", error);
  }
}

async function loadPetSelect() {
  try {
    const response = await fetch(`${API_BASE}/pets`);
    const pets = await response.json();
    const select = document.getElementById("pet-select");
    select.innerHTML =
      '<option value="">Выберите питомца</option>' +
      pets
        .map(
          (pet) =>
            `<option value="${pet.id}">${pet.name} (${pet.owner_name})</option>`
        )
        .join("");
  } catch (error) {
    console.error("Error loading pets for select:", error);
  }
}

//подробная информация о питомце
async function showPetDetail(petId) {
  try {
    const [petResponse, recordsResponse] = await Promise.all([
      fetch(`${API_BASE}/pets/${petId}`),
      fetch(`${API_BASE}/health-records/pet/${petId}`),
    ]);

    const pet = await petResponse.json();
    const records = await recordsResponse.json();

    displayPetDetail(pet, records);
    showHealthChart(records);

    document.querySelectorAll(".section").forEach((section) => {
      section.classList.remove("active");
    });

    document.getElementById("pet-detail-section").classList.add("active");
  } catch (error) {
    console.error("Error loading pet detail:", error);
  }
}

function displayPetDetail(pet, records) {
  const container = document.getElementById("pet-detail");
  container.innerHTML = `
        <div class="pet-detail-header">
            <div class="pet-info">
                <h2>${pet.name}</h2>
                <p><strong>Вид:</strong> ${pet.species}</p>
                <p><strong>Порода:</strong> ${pet.breed || "Не указана"}</p>
                <p><strong>Дата рождения:</strong> ${formatDate(
                  pet.date_of_birth
                )}</p>
                <p><strong>Возраст:</strong> ${calculateAge(
                  pet.date_of_birth
                )}</p>
                <p><strong>Окрас:</strong> ${pet.color || "Не указан"}</p>
                <p><strong>Микрочип:</strong> ${
                  pet.microchip_id || "Не установлен"
                }</p>
            </div>
            <a href="#" class="owner-link" onclick="showOwnerDetail(${
              pet.owner_id
            }); return false;">
                Владелец: ${pet.owner_name}
            </a>
        </div>
    `;

  const recordsContainer = document.getElementById("health-records");
  recordsContainer.innerHTML =
    `<h3>История визитов (${records.length})</h3>` +
    records
      .map(
        (record) => `
            <div class="health-record">
                <h4>${formatDate(record.visit_date)}</h4>
                <p><strong>Вес:</strong> ${record.weight || "Не указан"} кг</p>
                <p><strong>Температура:</strong> ${
                  record.temperature || "Не указана"
                }°C</p>
                <p><strong>Пульс:</strong> ${
                  record.heart_rate || "Не указан"
                } уд/мин</p>
                <p><strong>Дыхание:</strong> ${
                  record.respiratory_rate || "Не указано"
                } дых/мин</p>
                ${
                  record.diagnosis
                    ? `<p><strong>Диагноз:</strong> ${record.diagnosis}</p>`
                    : ""
                }
                ${
                  record.treatment
                    ? `<p><strong>Лечение:</strong> ${record.treatment}</p>`
                    : ""
                }
                ${
                  record.notes
                    ? `<p><strong>Заметки:</strong> ${record.notes}</p>`
                    : ""
                }
                ${
                  record.next_visit_date
                    ? `<p><strong>Следующий визит:</strong> ${formatDate(
                        record.next_visit_date
                      )}</p>`
                    : ""
                }
            </div>
        `
      )
      .join("");
}

//подробная информация о владельце
async function showOwnerDetail(ownerId) {
  try {
    const [ownerResponse, petsResponse] = await Promise.all([
      fetch(`${API_BASE}/owners/${ownerId}`),
      fetch(`${API_BASE}/pets/owner/${ownerId}`),
    ]);

    const owner = await ownerResponse.json();
    const pets = await petsResponse.json();

    displayOwnerDetail(owner, pets);

    document.querySelectorAll(".section").forEach((section) => {
      section.classList.remove("active");
    });

    document.getElementById("owner-detail-section").classList.add("active");
  } catch (error) {
    console.error("Error loading owner detail:", error);
  }
}

function displayOwnerDetail(owner, pets) {
  const container = document.getElementById("owner-detail");
  container.innerHTML = `
        <h2>${owner.first_name} ${owner.last_name}</h2>
        <p><strong>Email:</strong> ${owner.email}</p>
        <p><strong>Телефон:</strong> ${owner.phone || "Не указан"}</p>
        <p><strong>Адрес:</strong> ${owner.address || "Не указан"}</p>
        <p><strong>Зарегистрирован:</strong> ${formatDate(owner.created_at)}</p>
    `;

  const petsContainer = document.getElementById("owner-pets");
  petsContainer.innerHTML =
    `<h3>Питомцы (${pets.length})</h3>` +
    (pets.length > 0
      ? pets
          .map(
            (pet) => `
                <div class="card pet-card" onclick="showPetDetail(${pet.id})">
                    <h4>${pet.name}</h4>
                    <p><strong>Вид:</strong> ${pet.species}</p>
                    <p><strong>Порода:</strong> ${pet.breed || "Не указана"}</p>
                    <p><strong>Возраст:</strong> ${calculateAge(
                      pet.date_of_birth
                    )}</p>
                </div>
            `
          )
          .join("")
      : "<p>У этого владельца пока нет зарегистрированных питомцев.</p>");
}

function showHealthChart(records) {
  const ctx = document.getElementById("health-chart").getContext("2d");

  const sortedRecords = [...records].sort(
    (a, b) => new Date(a.visit_date) - new Date(b.visit_date)
  );

  const dates = sortedRecords.map((record) => formatDate(record.visit_date));

  const weights = sortedRecords
    .map((record) => record.weight)
    .filter((w) => w != null);

  const temperatures = sortedRecords
    .map((record) => record.temperature)
    .filter((t) => t != null);

  new Chart(ctx, {
    type: "line",
    data: {
      labels: dates,
      datasets: [
        {
          label: "Вес (кг)",
          data: weights,
          borderColor: "#3498db",
          backgroundColor: "rgba(52, 152, 219, 0.1)",
          yAxisID: "y",
        },
        {
          label: "Температура (°C)",
          data: temperatures,
          borderColor: "#e74c3c",
          backgroundColor: "rgba(231, 76, 60, 0.1)",
          yAxisID: "y1",
        },
      ],
    },
    options: {
      responsive: true,
      interaction: {
        mode: "index",
        intersect: false,
      },
      scales: {
        y: {
          type: "linear",
          display: true,
          position: "left",
          title: {
            display: true,
            text: "Вес (кг)",
          },
        },
        y1: {
          type: "linear",
          display: true,
          position: "right",
          title: {
            display: true,
            text: "Температура (°C)",
          },
          grid: {
            drawOnChartArea: false,
          },
        },
      },
    },
  });
}

async function handleAddPet(e) {
  e.preventDefault();
  const formData = new FormData(this);
  const petData = Object.fromEntries(formData);

  try {
    const response = await fetch(`${API_BASE}/pets`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(petData),
    });

    if (response.ok) {
      alert("Питомец успешно добавлен!");
      this.reset();
      loadPets();
      loadPetSelect();
      showSection("pets");
    } else {
      alert("Ошибка при добавлении питомца");
    }
  } catch (error) {
    console.error("Error adding pet:", error);
    alert("Ошибка при добавлении питомца");
  }
}

async function handleAddHealthRecord(e) {
  e.preventDefault();
  const formData = new FormData(this);
  const recordData = Object.fromEntries(formData);

  try {
    const response = await fetch(`${API_BASE}/health-records`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(recordData),
    });

    if (response.ok) {
      alert("Запись о здоровье успешно добавлена!");
      this.reset();
      showSection("pets");
    } else {
      alert("Ошибка при добавлении записи");
    }
  } catch (error) {
    console.error("Error adding health record:", error);
    alert("Ошибка при добавлении записи");
  }
}

async function handleAddOwner(e) {
  e.preventDefault();
  const formData = new FormData(this);
  const ownerData = Object.fromEntries(formData);

  if (!ownerData.first_name || !ownerData.last_name || !ownerData.email) {
    alert("Пожалуйста, заполните обязательные поля (Имя, Фамилия, Email)");
    return;
  }

  try {
    const response = await fetch(`${API_BASE}/owners`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(ownerData),
    });

    if (response.ok) {
      const newOwner = await response.json();
      alert(
        `Владелец ${newOwner.first_name} ${newOwner.last_name} успешно добавлен!`
      );
      this.reset();

      loadOwners();
      loadOwnerSelect();

      showSection("owners");
    } else {
      const error = await response.json();
      alert(`Ошибка: ${error.error || "Не удалось добавить владельца"}`);
    }
  } catch (error) {
    console.error("Error adding owner:", error);
    alert("Ошибка соединения с сервером");
  }
}

//вспомогательные функции
function calculateAge(birthDate) {
  if (!birthDate) return "Неизвестно";

  const birth = new Date(birthDate);
  const now = new Date();
  let years = now.getFullYear() - birth.getFullYear();
  let months = now.getMonth() - birth.getMonth();

  if (months < 0) {
    years--;
    months += 12;
  }

  if (years === 0) {
    return `${months} мес.`;
  } else {
    return `${years} г. ${months} мес.`;
  }
}

function formatDate(dateString) {
  if (!dateString) return "Не указана";

  const date = new Date(dateString);
  return date.toLocaleDateString("ru-RU");
}
