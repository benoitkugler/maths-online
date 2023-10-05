String formatTime(DateTime time) {
  time = time.toLocal();
  return "${_days[time.weekday]} ${time.day} ${_months[time.month]} ${time.year}, ${time.hour}h${time.minute.toString().padLeft(2, "0")}";
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
