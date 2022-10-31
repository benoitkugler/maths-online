String formatTime(DateTime time) {
  return "${_days[time.weekday]} ${time.day} ${_months[time.month]} ${time.year}, ${time.hour}h";
}

const _days = [
  "",
  "Lundi",
  "Mardi",
  "Mercredi",
  "Jeudi",
  "Vendredi",
  "Samedi",
  "Dimanche",
];

const _months = [
  "",
  "Jan.",
  "Fév.",
  "Mars",
  "Avril",
  "Mai",
  "Juin",
  "Juil.",
  "Août",
  "Sept.",
  "Oct.",
  "Nov.",
  "Déc.",
];
