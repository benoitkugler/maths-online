
	-- DO NOT EDIT - autogenerated by structgen 
		   
	
CREATE OR REPLACE FUNCTION f_delfunc (OUT func_dropped int
)
AS $func$
DECLARE
    _sql text;
BEGIN
    SELECT
        count(*)::int,
        'DROP FUNCTION ' || string_agg(oid::regprocedure::text, '; DROP FUNCTION ')
    FROM
        pg_proc
    WHERE
        starts_with (proname, 'structgen_validate_json')
        AND pg_function_is_visible(oid) INTO func_dropped,
        _sql;
    -- only returned if trailing DROPs succeed
    IF func_dropped > 0 THEN
        -- only if function(s) found
        EXECUTE _sql;
    END IF;
END
$func$
LANGUAGE plpgsql;

SELECT
    f_delfunc ();

DROP FUNCTION f_delfunc;


	CREATE OR REPLACE FUNCTION structgen_validate_json_string (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		RETURN jsonb_typeof(data) = 'string';
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_number (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		RETURN jsonb_typeof(data) = 'number';
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3340897463 (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		IF jsonb_typeof(data) != 'object' THEN 
			RETURN FALSE;
		END IF;
		RETURN (SELECT bool_and( 
			key IN ('expression', 'variable')
		) FROM jsonb_each(data))  
		AND structgen_validate_json_string(data->'expression')
AND structgen_validate_json_number(data->'variable')
		;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_3340897463 (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF; 
		RETURN (SELECT bool_and( structgen_validate_json_struct_3340897463(value) )  FROM jsonb_array_elements(data)) 
			;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

CREATE TABLE exercices (
	id serial PRIMARY KEY,
	title varchar  NOT NULL,
	description varchar  NOT NULL,
	random_parameters jsonb  CONSTRAINT random_parameters_structgen_validate_json_array_struct_3340897463 CHECK (structgen_validate_json_array_struct_3340897463(random_parameters))
);

	CREATE OR REPLACE FUNCTION structgen_validate_json_boolean (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		RETURN jsonb_typeof(data) = 'boolean';
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_struct_3651341631 (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		IF jsonb_typeof(data) != 'object' THEN 
			RETURN FALSE;
		END IF;
		RETURN (SELECT bool_and( 
			key IN ('Content', 'IsExpression')
		) FROM jsonb_each(data))  
		AND structgen_validate_json_string(data->'Content')
AND structgen_validate_json_boolean(data->'IsExpression')
		;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_struct_3651341631 (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF; 
		RETURN (SELECT bool_and( structgen_validate_json_struct_3651341631(value) )  FROM jsonb_array_elements(data)) 
			;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

CREATE TABLE formulas (
	Chunks jsonb  CONSTRAINT Chunks_structgen_validate_json_array_struct_3651341631 CHECK (structgen_validate_json_array_struct_3651341631(Chunks)),
	IsInline boolean  NOT NULL
);

CREATE TABLE formula_fields (
	Expression varchar  NOT NULL
);

CREATE TABLE list_fields (
	Choices varchar[] 
);

CREATE TABLE number_fields (
	Expression varchar  NOT NULL
);

	-- No validation : accept anything
	CREATE OR REPLACE FUNCTION structgen_validate_json_ (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		RETURN TRUE;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

	CREATE OR REPLACE FUNCTION structgen_validate_json_array_ (data jsonb)
		RETURNS boolean
		AS $f$
	BEGIN
		IF jsonb_typeof(data) = 'null' THEN RETURN TRUE; END IF;
		IF jsonb_typeof(data) != 'array' THEN RETURN FALSE; END IF;
		IF jsonb_array_length(data) = 0 THEN RETURN TRUE; END IF; 
		RETURN (SELECT bool_and( structgen_validate_json_(value) )  FROM jsonb_array_elements(data)) 
			;
	END;
	$f$
	LANGUAGE 'plpgsql'
	IMMUTABLE;

CREATE TABLE questions (
	id_exercice integer  NOT NULL,
	title varchar  NOT NULL,
	enonce jsonb  CONSTRAINT enonce_structgen_validate_json_array_ CHECK (structgen_validate_json_array_(enonce))
);

CREATE TABLE text_blocks (
	Text varchar  NOT NULL
);
ALTER TABLE questions ADD FOREIGN KEY(id_exercice) REFERENCES exercices ON DELETE CASCADE;