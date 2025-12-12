const API_URL = "https://talent-backend-986559184516.us-central1.run.app/talent";
const manualFieldIds = ["possibleRole", "tags", "notes"];

const dom = {
    form: document.getElementById("talentForm"),
    submitBtn: document.getElementById("submitBtn"),
    status: document.getElementById("status"),
    fullName: document.getElementById("fullName"),
    headline: document.getElementById("headline"),
    currentRole: document.getElementById("currentRole"),
    currentCompany: document.getElementById("currentCompany"),
    profileUrl: document.getElementById("profileUrl"),
    possibleRole: document.getElementById("possibleRole"),
    tags: document.getElementById("tags"),
    notes: document.getElementById("notes")
};

document.addEventListener("DOMContentLoaded", () => {
    bindEvents();
    restoreManualFields().finally(refreshLinkedInData);
});

function bindEvents() {
    dom.form.addEventListener("submit", handleSubmit);
    dom.form.addEventListener("input", handleManualFieldChange);
}

function restoreManualFields() {
    return new Promise((resolve) => {
        if (!chrome.storage || !chrome.storage.sync) {
            resolve();
            return;
        }
        chrome.storage.sync.get("manualDefaults", (stored) => {
            if (chrome.runtime?.lastError) {
                console.warn("Não foi possível restaurar campos manuais:", chrome.runtime.lastError);
                resolve();
                return;
            }
            const defaults = stored?.manualDefaults || {};
            manualFieldIds.forEach((id) => {
                if (defaults[id]) {
                    dom[id].value = defaults[id];
                }
            });
            resolve();
        });
    });
}

function handleManualFieldChange(event) {
    if (!manualFieldIds.includes(event.target.id)) {
        return;
    }
    persistManualFields();
}

function persistManualFields() {
    if (!chrome.storage || !chrome.storage.sync) {
        return;
    }
    const manualDefaults = {};
    manualFieldIds.forEach((id) => {
        manualDefaults[id] = dom[id].value || "";
    });
    chrome.storage.sync.set({ manualDefaults });
}

async function refreshLinkedInData(showFeedback = false) {
    if (showFeedback) {
        setStatus("Coletando dados do LinkedIn...", "info");
    }
    try {
        const tab = await getActiveTab();
        if (!tab) {
            setStatus("Não foi possível ler a aba atual.", "error");
            return;
        }
        const needsLinkedIn = !tab.url?.includes("linkedin.com");
        if (needsLinkedIn) {
            setStatus("Abra um perfil do LinkedIn para coletar os dados automaticamente.", "error");
            dom.profileUrl.value = tab.url || "";
            return;
        }

        const injectionResults = await chrome.scripting.executeScript({
            target: { tabId: tab.id },
            func: scrapeLinkedInProfile
        });
        const [result] = injectionResults || [];
        if (!result || !result.result) {
            setStatus("Não foi possível coletar os dados automaticamente.", "error");
            return;
        }
        fillFormWithLinkedIn(result.result);
        setStatus("Dados coletados do LinkedIn. Revise antes de enviar.", "success");
    } catch (error) {
        console.error("Erro ao coletar dados:", error);
        setStatus("Erro ao coletar dados. Preencha manualmente se necessário.", "error");
    }
}

function fillFormWithLinkedIn(data) {
    dom.fullName.value = data.fullName || dom.fullName.value;
    dom.headline.value = data.headline || dom.headline.value;
    dom.currentRole.value = data.currentRole || dom.currentRole.value;
    dom.currentCompany.value = data.currentCompany || dom.currentCompany.value;
    dom.profileUrl.value = data.profileUrl || dom.profileUrl.value;
}

async function handleSubmit(event) {
    event.preventDefault();
    const payload = buildPayload();
    if (!payload.full_name) {
        setStatus("Nome completo é obrigatório.", "error");
        dom.fullName.focus();
        return;
    }
    if (!payload.profile_url) {
        setStatus("Informe a URL do perfil.", "error");
        dom.profileUrl.focus();
        return;
    }
    try {
        toggleSubmitting(true);
        setStatus("Enviando dados para a API...", "info");
        const response = await fetch(API_URL, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(payload)
        });
        if (!response.ok) {
            const body = await safeReadJson(response);
            throw new Error(body?.message || `Erro ${response.status}`);
        }
        clearManualFields();
        setStatus("Talento registrado com sucesso!", "success");
    } catch (error) {
        console.error("Erro ao enviar talento:", error);
        setStatus(`Falha ao enviar: ${error.message}`, "error");
    } finally {
        toggleSubmitting(false);
    }
}

function buildPayload() {
    return {
        full_name: dom.fullName.value.trim(),
        headline: dom.headline.value.trim(),
        current_role: dom.currentRole.value.trim(),
        current_company: dom.currentCompany.value.trim(),
        profile_url: dom.profileUrl.value.trim(),
        possible_role: dom.possibleRole.value.trim(),
        tags: dom.tags.value
            .split(",")
            .map((tag) => tag.trim())
            .filter(Boolean),
        notes: dom.notes.value.trim()
    };
}

function toggleSubmitting(isSubmitting) {
    dom.submitBtn.disabled = isSubmitting;
}

function setStatus(message, type = "info") {
    dom.status.textContent = message;
    dom.status.className = `status ${type}`;
}

async function getActiveTab() {
    const tabs = await chrome.tabs.query({ active: true, currentWindow: true });
    return tabs && tabs.length ? tabs[0] : null;
}

async function safeReadJson(response) {
    try {
        return await response.json();
    } catch (error) {
        return null;
    }
}

function scrapeLinkedInProfile() {
    const sanitize = (value) => (value ? value.replace(/\s+/g, " ").trim() : "");
    const pickText = (selectors, root = document) => {
        for (const selector of selectors) {
            const el = root.querySelector(selector);
            if (el) {
                const text = sanitize(el.textContent);
                if (text) {
                    return text;
                }
            }
        }
        return "";
    };
    const parseRoleLine = (text) => {
        if (!text) {
            return { role: "", company: "" };
        }
        const cleaned = text.split(" · ")[0];
        const parts = cleaned.split(" at ");
        if (parts.length >= 2) {
            return {
                role: sanitize(parts[0]),
                company: sanitize(parts.slice(1).join(" at "))
            };
        }
        return { role: sanitize(cleaned), company: "" };
    };

    const fullName = pickText([
        ".pv-text-details__left-panel h1",
        "h1.text-heading-xlarge",
        ".ph5.pb5 h1",
        "main h1"
    ]);
    const headline = pickText([
        ".pv-text-details__left-panel .text-body-medium",
        ".pv-text-details__left-panel span.text-body-medium",
        ".text-body-medium.break-words",
        ".pv-top-card--list-bullet li",
        ".pv-text-details__right-panel span.text-body-small"
    ]);

    const roleLine = pickText([
        ".pv-text-details__right-panel li span[aria-hidden='true']",
        ".pv-text-details__right-panel li",
        ".pv-entity__summary-info h2 span[aria-hidden='true']",
        ".pv-entity__summary-info h3 span[aria-hidden='true']"
    ]);
    const parsedRole = parseRoleLine(roleLine);
    let currentRole = parsedRole.role;
    let currentCompany = parsedRole.company;

    const experienceSection =
        document.querySelector("#experience") ||
        document.querySelector("section[id*='experience']") ||
        document.querySelector("section[data-section='experience']") ||
        document.querySelector(".experience-section");
    if (experienceSection) {
        const directItems = Array.from(
            experienceSection.querySelectorAll(":scope ul > li, :scope .pvs-list > li")
        );
        const firstExperienceItem =
            directItems.find((item) => item.querySelector("[data-view-name='profile-component-entity']")) ||
            directItems[0] ||
            experienceSection;
        if (firstExperienceItem) {
            if (!currentRole) {
                currentRole = pickText(
                    [
                        ".mr1.hoverable-link-text.t-bold span[aria-hidden='true']",
                        ".display-flex.align-items-center.mr1.t-bold span[aria-hidden='true']",
                        ".t-bold span[aria-hidden='true']",
                        ".pvs-entity__summary-info h3 span[aria-hidden='true']",
                        "[data-field='experience_company_logo'] span[aria-hidden='true']"
                    ],
                    firstExperienceItem
                );
            }
            if (!currentCompany) {
                const companyLine = pickText(
                    [
                        ".t-14.t-normal span[aria-hidden='true']",
                        ".pvs-entity__secondary-title span[aria-hidden='true']",
                        ".display-flex.mt1 span[aria-hidden='true']",
                        ".pvs-entity__company-name span[aria-hidden='true']",
                        ".t-normal span[aria-hidden='true']"
                    ],
                    firstExperienceItem
                );
                currentCompany = companyLine ? sanitize(companyLine.split("·")[0]) : "";
            }
        }
    }

    return {
        fullName,
        headline,
        currentRole,
        currentCompany,
        profileUrl: window.location.href
    };
}

function clearManualFields() {
    manualFieldIds.forEach((id) => {
        dom[id].value = "";
    });
    persistManualFields();
}
