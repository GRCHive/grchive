package test.grchive.core.public.data.filters

import grchive.core.public.data.filters.EqualFilter
import io.kotest.matchers.shouldBe
import io.kotest.core.spec.style.StringSpec
import io.kotest.properties.assertAll

class EqualFilterTest: StringSpec({
    "EqualFilter<Byte>" {
        assertAll<Byte> { a ->
            EqualFilter<Byte>(a).target shouldBe a
        }
    }

    "EqualFilter<Short>" {
        assertAll<Short> { a ->
            EqualFilter<Short>(a).target shouldBe a
        }
    }

    "EqualFilter<Int>" {
        assertAll<Int> { a ->
            EqualFilter<Int>(a).target shouldBe a
        }
    }

    "EqualFilter<Long>" {
        assertAll<Long> { a ->
            EqualFilter<Long>(a).target shouldBe a
        }
    }

    "EqualFilter<Double>" {
        assertAll<Double> { a ->
            EqualFilter<Double>(a).target shouldBe a
        }
    }

    "EqualFilter<Float>" {
        assertAll<Float> { a ->
            EqualFilter<Float>(a).target shouldBe a
        }
    }

    "EqualFilter<Boolean>" {
        assertAll<Boolean> { a ->
            EqualFilter<Boolean>(a).target shouldBe a
        }
    }

    "EqualFilter<String>" {
        assertAll<String> { a ->
            EqualFilter<String>(a).target shouldBe a
        }
    }
})
