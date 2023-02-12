package questions

import (
	"fmt"
	"strings"

	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
)

const latexHeader = `
	% Custom commands and settings used by Isyro
	\newcommand{\isyroFieldHeight}{\phantom{$\sum_{\sum}^{\sum}$}} %% field height

	\newcommand{\isyroNumberField}{
		\fbox{\begin{minipage}{1cm}\isyroFieldHeight \end{minipage}}
	}
	
	\newcommand{\isyroExpressionField}[1]{
		\fbox{\begin{minipage}{#1 cm}\isyroFieldHeight \\ \isyroFieldHeight \end{minipage}}
	}
	
	\newcommand{\isyroQCMSquare}{\raisebox{-.25\height}{\huge$\square$}}
	
	\newcommand{\R}{\mathbb{R}}
	\newcommand{\Q}{\mathbb{Q}}
	\newcommand{\D}{\mathbb{D}}
	\newcommand{\Z}{\mathbb{Z}}
	\newcommand{\N}{\mathbb{N}}


`

func (qu QuestionInstance) toLatex(standalone bool) string {
	chunks := make([]string, len(qu.Enonce))
	// we add an extra new line between two text blocks
	isPreviousText := false
	for i, p := range qu.Enonce {
		_, isText := p.(TextInstance)

		if isPreviousText && isText {
			chunks[i] = "\n \\vspace{0.3cm} \n" + p.toLatex()
		} else {
			chunks[i] = p.toLatex()
		}

		isPreviousText = isText
	}
	// we add line breaks for clarity, they are ignored by the latex compiler anyway
	// we also disable indent inside one question
	code := fmt.Sprintf(`
	%% -------------- Question --------------
	{ \setlength{\parindent}{0pt}
	%s
	}`, strings.Join(chunks, "\n"))

	if standalone { // add custom defs
		return latexHeader + code
	}
	return code
}

func (qu QuestionInstance) ToLatex() string {
	return qu.toLatex(true)
}

func InstancesToLatex(questions []QuestionInstance) string {
	latexCodes := make([]string, len(questions))
	for i, qu := range questions {
		latexCodes[i] = `\item ` + qu.toLatex(false)
	}
	return fmt.Sprintf(`%s 
	
	\begin{enumerate}
	%s
	\end{enumerate}`, latexHeader, strings.Join(latexCodes, "\n"))
}

func lineToLatexCode(line client.TextLine) string {
	chunks := make([]string, len(line))
	for i, c := range line {
		if c.IsMath {
			chunks[i] = "$ " + c.Text + " $"
		} else {
			chunks[i] = c.Text
		}
	}
	out := strings.ReplaceAll(strings.Join(chunks, ""), "\n", `\\`)
	if strings.HasPrefix(out, `\\`) {
		out = `~` + out
	}
	return out
}

func (ti TextInstance) toLatex() string {
	text := lineToLatexCode(ti.Parts)

	attrs := ""
	if ti.Bold {
		attrs += `\bfseries `
	}
	if ti.Italic {
		attrs += `\itshape `
	}
	if ti.Smaller {
		attrs += `\small `
	}

	return fmt.Sprintf("{%s%s}", attrs, text)
}

func (fi FormulaDisplayInstance) toLatex() string {
	return "$$" + strings.Join(fi, " ") + "$$"
}

func (vi VariationTableInstance) toLatex() string { return "TODO" }
func (si SignTableInstance) toLatex() string      { return "TODO" }
func (fi FigureInstance) toLatex() string         { return "TODO" }
func (ti TableInstance) toLatex() string          { return "TODO" }
func (fi FunctionsGraphInstance) toLatex() string { return "TODO" }

func (ni NumberFieldInstance) toLatex() string { return `\isyroNumberField` }

func (ei ExpressionFieldInstance) toLatex() string {
	width := float64(ei.sizeHint())
	// map from 1 - 30 to 2cm - 10cm
	cm := 2. + (10-2)*(width-1)/(30-1)
	if ei.LabelLaTeX != "" { // add a new line and display mode
		return fmt.Sprintf(`$$ %s \isyroExpressionField{%.2f}$$`, ei.LabelLaTeX, cm)
	}
	return fmt.Sprintf(`\isyroExpressionField{%.2f}`, cm)
}

// requires the following latex packages
//   - \usepackage[inline]{enumitem}
//   - \usepackage{amssymb}
func (ri RadioFieldInstance) toLatex() string {
	props := ri.proposals()
	choices := make([]string, len(props))
	for i, p := range props {
		choices[i] = `\item ` + lineToLatexCode(p)
	}

	return fmt.Sprintf(`\begin{itemize}[label={\isyroQCMSquare}] %% vertical align
    %s
	\end{itemize}
	`, strings.Join(choices, "\n"))
}

func (di DropDownFieldInstance) toLatex() string {
	props := RadioFieldInstance(di).proposals()
	choices := make([]string, len(props))
	for i, p := range props {
		choices[i] = `\item ` + lineToLatexCode(p)
	}
	return fmt.Sprintf(`\begin{itemize*}[label={\isyroQCMSquare}] %% vertical align
    %s
	\end{itemize*}
	`, strings.Join(choices, "\n"))
}

func (oi OrderedListFieldInstance) toLatex() string {
	props := oi.proposals()
	choices := make([]string, len(props))
	for i, p := range props {
		choices[i] = lineToLatexCode(p)
	}
	box := fmt.Sprintf(`\isyroExpressionField{13}`)
	label := ""
	if oi.Label != "" {
		label = "$" + oi.Label + "$"
	}
	return fmt.Sprintf(`~\\ \begin{center} %s %s \\ 
	\vspace{0.5em}
	\textit{\small \'Eléments à ordonner:} %s
	\end{center}
	`, label, box, strings.Join(choices, " , "))
}

func (fi FigurePointFieldInstance) toLatex() string      { return "TODO" }
func (fi FigureVectorFieldInstance) toLatex() string     { return "TODO" }
func (vi VariationTableFieldInstance) toLatex() string   { return "TODO" }
func (si SignTableFieldInstance) toLatex() string        { return "TODO" }
func (fi FunctionPointsFieldInstance) toLatex() string   { return "TODO" }
func (fi FigureVectorPairFieldInstance) toLatex() string { return "TODO" }
func (fi FigureAffineLineFieldInstance) toLatex() string { return "TODO" }
func (ti TreeFieldInstance) toLatex() string             { return "TODO" }
func (pi ProofFieldInstance) toLatex() string            { return "TODO" }
func (pi TableFieldInstance) toLatex() string            { return "TODO" }
func (pi VectorFieldInstance) toLatex() string           { return "TODO" }
