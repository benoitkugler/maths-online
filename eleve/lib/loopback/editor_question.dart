/// [QuestionPage] is the editor version of a question,
/// and is not used by the loopback app itself,
/// so that is treated as an opaque type
typedef QuestionPage = dynamic;

QuestionPage questionPageFromJson(dynamic json) {
  return json;
}

dynamic questionPageToJson(QuestionPage qp) {
  return qp;
}
